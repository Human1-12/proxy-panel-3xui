package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/mhsanaei/3x-ui/v3/internal/database/model"
)

// OneClickRealityRequest carries the parameters for one-click batch generation
// of VLESS + TCP + REALITY + Vision inbounds. All fields are optional; empty
// values fall back to sensible defaults (10 nodes starting at port 20000).
type OneClickRealityRequest struct {
	Count        int      `json:"count"`
	PortStart    int      `json:"portStart"`
	RemarkPrefix string   `json:"remarkPrefix"`
	Dest         string   `json:"dest"`
	ServerNames  []string `json:"serverNames"`
}

// OneClickCreatedInbound summarizes one successfully created inbound.
type OneClickCreatedInbound struct {
	Id     int    `json:"id"`
	Port   int    `json:"port"`
	Remark string `json:"remark"`
	Email  string `json:"email"`
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

// BatchCreateRealityVision generates req.Count VLESS + TCP + REALITY + Vision
// inbounds in one call. Each node gets its own freshly generated X25519 keypair,
// random shortId, client UUID and subId. Inbounds are persisted through the
// normal AddInbound path so they are identical to manually created ones.
// Returns the per-inbound result and whether xray needs a restart.
func (s *InboundService) BatchCreateRealityVision(userId int, req OneClickRealityRequest) (*OneClickResult, bool, error) {
	if req.Count <= 0 {
		req.Count = 10
	}
	if req.Count > 100 {
		req.Count = 100
	}
	if req.PortStart <= 0 {
		req.PortStart = 20000
	}
	if strings.TrimSpace(req.RemarkPrefix) == "" {
		req.RemarkPrefix = "reality"
	}
	if strings.TrimSpace(req.Dest) == "" {
		req.Dest = "www.microsoft.com:443"
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
	port := req.PortStart

	for i := 0; i < req.Count; i++ {
		kpAny, err := server.GetNewX25519Cert()
		if err != nil {
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

		clientUUID := uuid.New().String()
		email := fmt.Sprintf("%s-%d-%s", req.RemarkPrefix, i+1, oneClickRandHex(3))
		subId := oneClickRandHex(8)
		shortId := oneClickRandHex(4)

		settings := map[string]any{
			"clients": []map[string]any{{
				"id":         clientUUID,
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
				"shortIds":    []string{shortId},
				"settings": map[string]any{
					"publicKey":   pub,
					"fingerprint": "chrome",
					"spiderX":     "/",
				},
			},
			"tcpSettings": map[string]any{"header": map[string]any{"type": "none"}},
		}
		sniffing := map[string]any{
			"enabled":      false,
			"destOverride": []string{"http", "tls", "quic", "fakedns"},
			"metadataOnly": false,
			"routeOnly":    false,
		}

		sBytes, _ := json.Marshal(settings)
		stBytes, _ := json.Marshal(stream)
		snBytes, _ := json.Marshal(sniffing)

		inbound := &model.Inbound{
			UserId:         userId,
			Enable:         true,
			Protocol:       model.VLESS,
			Remark:         fmt.Sprintf("%s-%02d", req.RemarkPrefix, i+1),
			Listen:         "",
			Total:          0,
			ExpiryTime:     0,
			TrafficReset:   "never",
			Settings:       string(sBytes),
			StreamSettings: string(stBytes),
			Sniffing:       string(snBytes),
		}

		var created *model.Inbound
		var restart bool
		var addErr error
		placed := false
		for attempts := 0; attempts < 500; attempts++ {
			inbound.Port = port
			port++
			created, restart, addErr = s.AddInbound(inbound)
			if addErr == nil {
				placed = true
				break
			}
			// A port conflict just means "try the next port"; anything else is fatal for this node.
			if strings.Contains(strings.ToLower(addErr.Error()), "port") {
				continue
			}
			break
		}
		if !placed {
			result.Failed++
			msg := "unknown error"
			if addErr != nil {
				msg = addErr.Error()
			}
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", email, msg))
			continue
		}

		anyRestart = anyRestart || restart
		result.Created++
		result.Inbounds = append(result.Inbounds, OneClickCreatedInbound{
			Id:     created.Id,
			Port:   created.Port,
			Remark: created.Remark,
			Email:  email,
		})
	}

	return result, anyRestart, nil
}
