package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mhsanaei/3x-ui/v3/internal/logger"
)

// firewallBackend identifies which managed host-firewall manager is active.
type firewallBackend int

const (
	fwNone firewallBackend = iota
	fwUFW
	fwFirewalld
)

func (b firewallBackend) String() string {
	switch b {
	case fwUFW:
		return "ufw"
	case fwFirewalld:
		return "firewalld"
	default:
		return "none"
	}
}

// autoFirewallEnabled reports whether the panel should try to open host-firewall
// ports for newly created inbounds. Enabled unless XUI_AUTO_FIREWALL is set to a
// false-y value — mirrors the XUI_ENABLE_FAIL2BAN convention.
func autoFirewallEnabled() bool {
	v, ok := os.LookupEnv("XUI_AUTO_FIREWALL")
	if !ok {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "false", "0", "off", "no":
		return false
	default:
		return true
	}
}

// runFirewallCmd runs a firewall command with a short timeout and returns its
// combined output plus error, so every call site can log uniformly.
func runFirewallCmd(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var buf bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

// detectFirewallBackend picks the active managed host firewall in priority order:
// ufw (if active) → firewalld (if running); fwNone otherwise. We deliberately do
// NOT touch raw iptables/nftables: rules added there don't persist across reboot
// and can clash with other tooling, and a host with no managed firewall usually
// relies on a cloud security group the panel can't reach anyway.
func detectFirewallBackend() firewallBackend {
	// ufw prints "Status: active" when enabled.
	if out, err := runFirewallCmd("ufw", "status"); err == nil {
		if strings.Contains(strings.ToLower(out), "status: active") {
			return fwUFW
		}
	}
	// firewalld: `firewall-cmd --state` prints "running" and exits 0 when up.
	if out, err := runFirewallCmd("firewall-cmd", "--state"); err == nil {
		if strings.Contains(strings.ToLower(out), "running") {
			return fwFirewalld
		}
	}
	return fwNone
}

// openFirewallPorts opens the given TCP ports (and UDP too when alsoUDP is set)
// in whatever managed host firewall is active. Best-effort and idempotent: any
// failure is logged and folded into the returned human-readable summary, never
// propagated — a firewall problem must never fail node creation. The summary is
// surfaced in the one-click API response so the operator sees what happened.
func openFirewallPorts(ports []int, alsoUDP bool) string {
	if len(ports) == 0 {
		return ""
	}
	if !autoFirewallEnabled() {
		return "auto-firewall disabled (XUI_AUTO_FIREWALL); opened no ports"
	}
	backend := detectFirewallBackend()
	if backend == fwNone {
		msg := "no managed firewall (ufw/firewalld) active; skipped — ensure your cloud security group allows the new ports"
		logger.Info("one-click firewall:", msg)
		return msg
	}

	protos := []string{"tcp"}
	if alsoUDP {
		protos = append(protos, "udp")
	}

	opened := 0
	var failed []string
	for _, port := range ports {
		for _, proto := range protos {
			if err := openOnePort(backend, port, proto); err != nil {
				failed = append(failed, fmt.Sprintf("%d/%s", port, proto))
				logger.Warning("one-click firewall: failed to open", port, proto, "via", backend.String(), ":", err)
			} else {
				opened++
			}
		}
	}
	// firewalld applies --permanent rules only after a reload.
	if backend == fwFirewalld {
		if out, err := runFirewallCmd("firewall-cmd", "--reload"); err != nil {
			logger.Warning("one-click firewall: firewalld reload failed:", out, err)
		}
	}

	summary := fmt.Sprintf("%s: opened %d rule(s) across %d port(s)", backend.String(), opened, len(ports))
	if len(failed) > 0 {
		summary += "; failed: " + strings.Join(failed, ", ")
	}
	logger.Info("one-click firewall:", summary)
	return summary
}

// openOnePort adds a single idempotent allow rule for one port/proto.
func openOnePort(backend firewallBackend, port int, proto string) error {
	p := strconv.Itoa(port)
	switch backend {
	case fwUFW:
		// `ufw allow <port>/<proto>` is idempotent (skips an existing rule).
		_, err := runFirewallCmd("ufw", "allow", p+"/"+proto)
		return err
	case fwFirewalld:
		// --permanent persists across reboot; the caller reloads afterwards to
		// apply it to the running firewall. firewall-cmd reports ALREADY_ENABLED
		// as success (exit 0), so this is idempotent.
		_, err := runFirewallCmd("firewall-cmd", "--permanent", "--add-port="+p+"/"+proto)
		return err
	default:
		return fmt.Errorf("no firewall backend")
	}
}
