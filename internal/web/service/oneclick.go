package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/mhsanaei/3x-ui/v3/internal/database"
	"github.com/mhsanaei/3x-ui/v3/internal/database/model"
	"github.com/mhsanaei/3x-ui/v3/internal/util/random"
)

// OneClickRealityRequest carries the parameters for one-click batch generation
// of inbounds. All fields are optional; empty values fall back to sensible
// defaults (10 nodes starting at port 20000).
//
// Protocol selects the preset:
//   - "reality" (default): VLESS + TCP + REALITY + Vision (uses Dest/ServerNames)
//   - "ss2022":            Shadowsocks 2022-blake3-aes-256-gcm (Dest/ServerNames ignored)
//   - "vmess":             VMess + TCP, no TLS (Dest/ServerNames ignored)
//   - "vlessTcp":          VLESS + TCP, security "none", no REALITY (Dest/ServerNames ignored)
//
// All presets are certificate-free — they run on a bare VPS IP with no domain.
type OneClickRealityRequest struct {
	Count        int      `json:"count"`
	PortStart    int      `json:"portStart"`
	RemarkPrefix string   `json:"remarkPrefix"`
	Protocol     string   `json:"protocol"`
	Dest         string   `json:"dest"`
	ServerNames  []string `json:"serverNames"`
}

// OneClickCreatedInbound summarizes one successfully created inbound.
type OneClickCreatedInbound struct {
	Id       int    `json:"id"`
	Port     int    `json:"port"`
	Remark   string `json:"remark"`
	Email    string `json:"email"`
	Protocol string `json:"protocol"`
}

// OneClickResult is the outcome of a batch generation run.
type OneClickResult struct {
	Requested int                      `json:"requested"`
	Created   int                      `json:"created"`
	Failed    int                      `json:"failed"`
	Inbounds  []OneClickCreatedInbound `json:"inbounds"`
	Errors    []string                 `json:"errors,omitempty"`
}

// oneClickRandHex returns nBytes of cryptographically-random data as a lowercase
// hex string (2*nBytes characters).
func oneClickRandHex(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand should never fail; degrade to a uuid-derived value.
		return strings.ReplaceAll(uuid.New().String(), "-", "")[:nBytes*2]
	}
	return hex.EncodeToString(b)
}

// oneClickSniffing returns the default (disabled) sniffing block shared by the presets.
func oneClickSniffing() map[string]any {
	return map[string]any{
		"enabled":      false,
		"destOverride": []string{"http", "tls", "quic", "fakedns"},
		"metadataOnly": false,
		"routeOnly":    false,
	}
}

// realityInboundParts builds the settings / streamSettings / sniffing maps for a
// VLESS + TCP + REALITY + Vision node, given a freshly generated X25519 keypair.
func realityInboundParts(req OneClickRealityRequest, email, subId, priv, pub string) (map[string]any, map[string]any, map[string]any) {
	settings := map[string]any{
		"clients": []map[string]any{{
			"id":         uuid.New().String(),
			"flow":       "xtls-rprx-vision",
			"email":      email,
			"enable":     true,
			"subId":      subId,
			"reset":      0,
			"limitIp":    0,
			"totalGB":    0,
			"expiryTime": 0,
		}},
		"decryption": "none",
		"fallbacks":  []any{},
	}
	stream := map[string]any{
		"network":  "tcp",
		"security": "reality",
		"realitySettings": map[string]any{
			"show":        false,
			"dest":        req.Dest,
			"serverNames": req.ServerNames,
			"privateKey":  priv,
			"shortIds":    []string{oneClickRandHex(4)},
			"settings": map[string]any{
				"publicKey":   pub,
				"fingerprint": "chrome",
				"spiderX":     "/",
			},
		},
		"tcpSettings": map[string]any{"header": map[string]any{"type": "none"}},
	}
	return settings, stream, oneClickSniffing()
}

// ss2022InboundParts builds the settings / streamSettings / sniffing maps for a
// Shadowsocks 2022-blake3-aes-256-gcm node. That method needs a 32-byte PSK, so
// both the server key and each client key are 32 random bytes, base64-encoded —
// matching what the panel produces for a hand-created SS-2022 inbound (and what
// AddInbound's normalizeShadowsocksClientKeys validates).
func ss2022InboundParts(email, subId string) (map[string]any, map[string]any, map[string]any) {
	const ssKeyBytes = 32 // 2022-blake3-aes-256-gcm
	settings := map[string]any{
		"method":   "2022-blake3-aes-256-gcm",
		"password": random.Base64Bytes(ssKeyBytes),
		"network":  "tcp,udp",
		"clients": []map[string]any{{
			"method":     "",
			"password":   random.Base64Bytes(ssKeyBytes),
			"email":      email,
			"enable":     true,
			"subId":      subId,
			"reset":      0,
			"limitIp":    0,
			"totalGB":    0,
			"expiryTime": 0,
			"comment":    "",
		}},
		"ivCheck": false,
	}
	stream := map[string]any{
		"network":     "tcp",
		"security":    "none",
		"tcpSettings": map[string]any{"header": map[string]any{"type": "none"}},
	}
	return settings, stream, oneClickSniffing()
}

// vmessTcpInboundParts builds the settings / streamSettings / sniffing maps for a
// VMess + TCP node with no TLS (certificate-free). Only a client UUID is needed —
// no keygen. Deliberately omits alterId (the panel uses VMessAEAD; alterId would
// be legacy). The client's security:"auto" is stripped before xray sees it, by
// design — same as a hand-created VMess inbound.
func vmessTcpInboundParts(email, subId string) (map[string]any, map[string]any, map[string]any) {
	settings := map[string]any{
		"clients": []map[string]any{{
			"id":         uuid.New().String(),
			"security":   "auto",
			"email":      email,
			"enable":     true,
			"subId":      subId,
			"reset":      0,
			"limitIp":    0,
			"totalGB":    0,
			"expiryTime": 0,
			"tgId":       0,
			"comment":    "",
		}},
	}
	stream := map[string]any{
		"network":     "tcp",
		"security":    "none",
		"tcpSettings": map[string]any{"header": map[string]any{"type": "none"}},
	}
	return settings, stream, oneClickSniffing()
}

// vlessTcpInboundParts builds the settings / streamSettings / sniffing maps for
// the simplest VLESS + TCP node: security "none", no REALITY, no TLS. flow MUST
// stay empty — xtls-rprx-vision needs a TLS/REALITY transport and xray refuses to
// start with a vision flow under security:"none". encryption is omitted (it gets
// stripped before xray sees it anyway), matching the reality builder.
func vlessTcpInboundParts(email, subId string) (map[string]any, map[string]any, map[string]any) {
	settings := map[string]any{
		"clients": []map[string]any{{
			"id":         uuid.New().String(),
			"flow":       "",
			"email":      email,
			"enable":     true,
			"subId":      subId,
			"reset":      0,
			"limitIp":    0,
			"totalGB":    0,
			"expiryTime": 0,
		}},
		"decryption": "none",
		"fallbacks":  []any{},
	}
	stream := map[string]any{
		"network":     "tcp",
		"security":    "none",
		"tcpSettings": map[string]any{"header": map[string]any{"type": "none"}},
	}
	return settings, stream, oneClickSniffing()
}

// BatchCreateRealityVision generates req.Count inbounds in one call, one per
// node, each with its own freshly generated keys / UUID / subId. Ports are
// allocated deterministically from a pre-loaded set of used ports. Inbounds are
// persisted through the normal AddInbound path so they are identical to manually
// created ones. Returns the per-inbound result and whether xray needs a restart.
//
// Despite the historical name, req.Protocol selects the preset (see
// OneClickRealityRequest); "reality" is the default.
func (s *InboundService) BatchCreateRealityVision(userId int, req OneClickRealityRequest) (*OneClickResult, bool, error) {
	if req.Count <= 0 {
		req.Count = 10
	}
	if req.Count > 100 {
		req.Count = 100
	}
	if req.PortStart <= 0 || req.PortStart > 65535 {
		req.PortStart = 20000
	}
	req.Protocol = strings.ToLower(strings.TrimSpace(req.Protocol))
	// Known one-click presets → default remark prefix. Adding a preset means one
	// entry here + one case in the switch below (+ its builder). Unknown or empty
	// protocols fall back to reality via this whitelist, never silently accepted.
	// NOTE: keys are lowercase because req.Protocol is lowercased above, so the
	// frontend value "vlessTcp" arrives here as "vlesstcp".
	oneClickDefaultPrefix := map[string]string{
		"reality":  "reality",
		"ss2022":   "ss",
		"vmess":    "vmess",
		"vlesstcp": "vless",
	}
	if _, ok := oneClickDefaultPrefix[req.Protocol]; !ok {
		req.Protocol = "reality"
	}
	if strings.TrimSpace(req.RemarkPrefix) == "" {
		req.RemarkPrefix = oneClickDefaultPrefix[req.Protocol]
	}
	req.Dest = strings.TrimSpace(req.Dest)
	if req.Dest == "" {
		req.Dest = "www.microsoft.com:443"
	}
	// REALITY dest must be host:port; default to :443 when the caller omits the port.
	if !strings.Contains(req.Dest, ":") {
		req.Dest = req.Dest + ":443"
	}
	if len(req.ServerNames) == 0 {
		host := req.Dest
		if i := strings.LastIndex(host, ":"); i > 0 {
			host = host[:i]
		}
		req.ServerNames = []string{host}
	}

	server := &ServerService{}
	result := &OneClickResult{Requested: req.Count}
	anyRestart := false

	// Pre-load every port already taken (existing inbounds + the internal Xray
	// API port) so free ports can be handed out deterministically, instead of
	// probing AddInbound and pattern-matching on its error text.
	usedPorts := make(map[int]bool)
	var existingPorts []int
	if err := database.GetDB().Model(model.Inbound{}).Pluck("port", &existingPorts).Error; err == nil {
		for _, p := range existingPorts {
			usedPorts[p] = true
		}
	}
	usedPorts[reservedAPIPort()] = true

	nextPort := req.PortStart
	allocPort := func() (int, bool) {
		for nextPort <= 65535 {
			p := nextPort
			nextPort++
			if !usedPorts[p] {
				usedPorts[p] = true
				return p, true
			}
		}
		return 0, false
	}

	for i := 0; i < req.Count; i++ {
		email := fmt.Sprintf("%s-%d-%s", req.RemarkPrefix, i+1, oneClickRandHex(3))
		subId := oneClickRandHex(8)

		var proto model.Protocol
		var settings, stream, sniffing map[string]any

		switch req.Protocol {
		case "ss2022":
			proto = model.Shadowsocks
			settings, stream, sniffing = ss2022InboundParts(email, subId)
		case "vmess":
			proto = model.VMESS
			settings, stream, sniffing = vmessTcpInboundParts(email, subId)
		case "vlesstcp":
			proto = model.VLESS
			settings, stream, sniffing = vlessTcpInboundParts(email, subId)
		default: // reality
			proto = model.VLESS
			kpAny, err := server.GetNewX25519Cert()
			if err != nil {
				// GetNewX25519Cert shells out to the xray binary; if the very first
				// call fails (e.g. xray missing) every node would fail identically,
				// so fail the whole batch fast with a clear, actionable message.
				if i == 0 {
					return nil, false, fmt.Errorf("reality keygen failed (is the xray binary present?): %w", err)
				}
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("node %d: reality keygen failed: %v", i+1, err))
				continue
			}
			kp, _ := kpAny.(map[string]any)
			priv, _ := kp["privateKey"].(string)
			pub, _ := kp["publicKey"].(string)
			if priv == "" || pub == "" {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("node %d: empty reality key", i+1))
				continue
			}
			settings, stream, sniffing = realityInboundParts(req, email, subId, priv, pub)
		}

		sBytes, _ := json.Marshal(settings)
		stBytes, _ := json.Marshal(stream)
		snBytes, _ := json.Marshal(sniffing)

		inbound := &model.Inbound{
			UserId:         userId,
			Enable:         true,
			Protocol:       proto,
			Remark:         fmt.Sprintf("%s-%02d", req.RemarkPrefix, i+1),
			Listen:         "",
			Total:          0,
			ExpiryTime:     0,
			TrafficReset:   "never",
			Settings:       string(sBytes),
			StreamSettings: string(stBytes),
			Sniffing:       string(snBytes),
		}

		port, ok := allocPort()
		if !ok {
			// No free ports left in range; remaining nodes can't be placed either.
			remaining := req.Count - i
			result.Failed += remaining
			result.Errors = append(result.Errors, fmt.Sprintf(
				"ran out of free ports at/above %d; %d node(s) not created", req.PortStart, remaining))
			break
		}
		inbound.Port = port

		created, restart, addErr := s.AddInbound(inbound)
		if addErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s (port %d): %s", email, port, addErr.Error()))
			continue
		}

		anyRestart = anyRestart || restart
		result.Created++
		result.Inbounds = append(result.Inbounds, OneClickCreatedInbound{
			Id:       created.Id,
			Port:     created.Port,
			Remark:   created.Remark,
			Email:    email,
			Protocol: string(proto),
		})
	}

	return result, anyRestart, nil
}
