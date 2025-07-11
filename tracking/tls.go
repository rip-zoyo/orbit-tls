package tracking

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rip-zoyo/orbit-tls/fingerprint"
)

type TLSTracker struct {
	mu          sync.RWMutex
	connections map[string]*ConnectionDetails
}

type ConnectionDetails struct {
	ClientRandom         []byte    `json:"client_random"`
	SessionID            []byte    `json:"session_id"`
	CipherSuite          uint16    `json:"cipher_suite"`
	CompressionMethod    uint8     `json:"compression_method"`
	TLSVersion           uint16    `json:"tls_version"`
	Extensions           []uint16  `json:"extensions"`
	SupportedVersions    []uint16  `json:"supported_versions"`
	SupportedGroups      []uint16  `json:"supported_groups"`
	SignatureAlgorithms  []uint16  `json:"signature_algorithms"`
	ALPNProtocols        []string  `json:"alpn_protocols"`
	ServerName           string    `json:"server_name"`
	KeyShare             []byte    `json:"key_share"`
	PeerCertificates     [][]byte  `json:"peer_certificates"`
	HandshakeComplete    bool      `json:"handshake_complete"`
	ConnectedAt          time.Time `json:"connected_at"`
	HTTP2Settings        map[string]uint32 `json:"http2_settings"`
	HTTP2Frames          []fingerprint.Frame `json:"http2_frames"`
	HTTP2WindowUpdate    uint32           `json:"http2_window_update"`
	HTTP2Priority        *fingerprint.HeaderPriority `json:"http2_priority"`
}

type TrackedDialer struct {
	dialer  *net.Dialer
	tracker *TLSTracker
}

var GlobalTracker = NewTLSTracker()

func NewTLSTracker() *TLSTracker {
	return &TLSTracker{
		connections: make(map[string]*ConnectionDetails),
	}
}

func NewTrackedDialer() *TrackedDialer {
	return &TrackedDialer{
		dialer: &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
		tracker: GlobalTracker,
	}
}

func (td *TrackedDialer) DialTLS(network, addr string, config *tls.Config) (net.Conn, error) {
	details := &ConnectionDetails{
		ConnectedAt: time.Now(),
		ServerName:  config.ServerName,
	}

	trackedConfig := td.createTrackedTLSConfig(config, details)
	
	conn, err := tls.Dial(network, addr, trackedConfig)
	if err != nil {
		return nil, err
	}

	td.updateConnectionDetails(addr, conn, details)
	
	return conn, nil
}

func (td *TrackedDialer) createTrackedTLSConfig(original *tls.Config, details *ConnectionDetails) *tls.Config {
	config := original.Clone()
	
	config.VerifyConnection = func(cs tls.ConnectionState) error {
		details.TLSVersion = cs.Version
		details.CipherSuite = cs.CipherSuite
		
		if len(cs.PeerCertificates) > 0 {
			details.PeerCertificates = make([][]byte, len(cs.PeerCertificates))
			for i, cert := range cs.PeerCertificates {
				details.PeerCertificates[i] = cert.Raw
			}
		}
		
		details.HandshakeComplete = true
		
		if original.VerifyConnection != nil {
			return original.VerifyConnection(cs)
		}
		
		return nil
	}
	
	return config
}

func (td *TrackedDialer) updateConnectionDetails(addr string, conn *tls.Conn, details *ConnectionDetails) {
	state := conn.ConnectionState()
	
	details.TLSVersion = state.Version
	details.CipherSuite = state.CipherSuite
	details.HandshakeComplete = state.HandshakeComplete
	
	if len(state.PeerCertificates) > 0 {
		details.PeerCertificates = make([][]byte, len(state.PeerCertificates))
		for i, cert := range state.PeerCertificates {
			details.PeerCertificates[i] = cert.Raw
		}
	}
	
	td.tracker.StoreConnection(addr, details)
}

func (t *TLSTracker) StoreConnection(addr string, details *ConnectionDetails) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.connections[addr] = details
}

func (t *TLSTracker) GetConnection(addr string) (*ConnectionDetails, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	details, exists := t.connections[addr]
	return details, exists
}

func (t *TLSTracker) GetAllConnections() map[string]*ConnectionDetails {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	result := make(map[string]*ConnectionDetails)
	for k, v := range t.connections {
		result[k] = v
	}
	return result
}

func (t *TLSTracker) ClearConnections() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.connections = make(map[string]*ConnectionDetails)
}

func GenerateFingerprintData(profile interface{}, details *ConnectionDetails) *fingerprint.Data {
	if details == nil {
		return generateFallbackFingerprint(profile)
	}
	
	p := profile.(interface{ GetJA3() string })
	tlsVersion, cipherSuites, extensions, supportedGroups, _, err := fingerprint.ParseJA3(p.GetJA3())
	if err != nil {
		return generateFallbackFingerprint(profile)
	}
	
	fp := &fingerprint.Data{
		TLSVersion:           fmt.Sprintf("%d", details.TLSVersion),
		TLSVersionRecord:     fmt.Sprintf("%d", tlsVersion),
		TLSVersionNegotiated: fmt.Sprintf("%d", details.TLSVersion),
		JA3:                  p.GetJA3(),
		JA3Hash:              fingerprint.GenerateJA3Hash(p.GetJA3()),
		ClientRandom:         hex.EncodeToString(details.ClientRandom),
		SessionID:            hex.EncodeToString(details.SessionID),
	}
	
	fp.JA4 = fingerprint.GenerateJA4(details.TLSVersion, cipherSuites, extensions, supportedGroups, details.ALPNProtocols)
	fp.JA4_R = fingerprint.GenerateJA4R(details.TLSVersion, cipherSuites, extensions, supportedGroups, details.ALPNProtocols)
	fp.PeetPrint = fingerprint.GeneratePeetPrint(details.TLSVersion, tlsVersion, cipherSuites, extensions, supportedGroups, details.SignatureAlgorithms)
	fp.PeetPrintHash = fingerprint.GeneratePeetPrintHash(fp.PeetPrint)
	
	fp.CipherSuites = make([]string, len(cipherSuites))
	for i, cipher := range cipherSuites {
		fp.CipherSuites[i] = fingerprint.GetCipherSuiteName(cipher)
	}
	
	fp.Extensions = make([]fingerprint.Extension, len(extensions))
	for i, ext := range extensions {
		fp.Extensions[i] = fingerprint.Extension{
			Type: ext,
			Name: fingerprint.GetExtensionName(ext),
		}
	}
	
	fp.SupportedGroups = make([]string, len(supportedGroups))
	for i, group := range supportedGroups {
		fp.SupportedGroups[i] = getSupportedGroupName(group)
	}
	
	fp.SignatureAlgorithms = make([]string, len(details.SignatureAlgorithms))
	for i, alg := range details.SignatureAlgorithms {
		fp.SignatureAlgorithms[i] = getSignatureAlgorithmName(alg)
	}
	
	if len(details.HTTP2Settings) > 0 {
		fp.HTTP2 = &fingerprint.HTTP2Data{
			Settings:         details.HTTP2Settings,
			WindowUpdate:     details.HTTP2WindowUpdate,
			HeaderPriority:   details.HTTP2Priority,
			SentFrames:       details.HTTP2Frames,
			ConnectionPreface: "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n",
		}
		fp.AkamaiFP = generateAkamaiFingerprint(details.HTTP2Settings, details.HTTP2WindowUpdate, details.HTTP2Priority)
		fp.AkamaiFPHash = fingerprint.GenerateJA3Hash(fp.AkamaiFP)
	}
	
	return fp
}

func generateFallbackFingerprint(profile interface{}) *fingerprint.Data {
	p := profile.(interface{ GetJA3() string })
	return &fingerprint.Data{
		JA3:           p.GetJA3(),
		JA3Hash:       fingerprint.GenerateJA3Hash(p.GetJA3()),
		ClientRandom:  generateRandomHex(32),
		SessionID:     generateRandomHex(32),
	}
}

func generateRandomHex(length int) string {
	return hex.EncodeToString(make([]byte, length))
}

func getSupportedGroupName(group uint16) string {
	groups := map[uint16]string{
		0x001d: "x25519",
		0x0017: "secp256r1",
		0x0018: "secp384r1",
		0x0019: "secp521r1",
		0x001e: "x448",
	}
	if name, exists := groups[group]; exists {
		return name
	}
	return fmt.Sprintf("UNKNOWN_GROUP_%d", group)
}

func getSignatureAlgorithmName(alg uint16) string {
	algs := map[uint16]string{
		0x0401: "rsa_pkcs1_sha256",
		0x0501: "rsa_pkcs1_sha384",
		0x0601: "rsa_pkcs1_sha512",
		0x0403: "ecdsa_secp256r1_sha256",
		0x0503: "ecdsa_secp384r1_sha384",
		0x0603: "ecdsa_secp521r1_sha512",
		0x0804: "rsa_pss_rsae_sha256",
		0x0805: "rsa_pss_rsae_sha384",
		0x0806: "rsa_pss_rsae_sha512",
		0x0807: "ed25519",
		0x0808: "ed448",
	}
	if name, exists := algs[alg]; exists {
		return name
	}
	return fmt.Sprintf("UNKNOWN_SIG_ALG_%d", alg)
}

func generateAkamaiFingerprint(settings map[string]uint32, windowUpdate uint32, priority *fingerprint.HeaderPriority) string {
	var parts []string
	
	if val, ok := settings["HEADER_TABLE_SIZE"]; ok {
		parts = append(parts, fmt.Sprintf("1:%d", val))
	}
	if val, ok := settings["ENABLE_PUSH"]; ok {
		parts = append(parts, fmt.Sprintf("2:%d", val))
	}
	if val, ok := settings["MAX_CONCURRENT_STREAMS"]; ok {
		parts = append(parts, fmt.Sprintf("3:%d", val))
	}
	if val, ok := settings["INITIAL_WINDOW_SIZE"]; ok {
		parts = append(parts, fmt.Sprintf("4:%d", val))
	}
	if val, ok := settings["MAX_FRAME_SIZE"]; ok {
		parts = append(parts, fmt.Sprintf("5:%d", val))
	}
	if val, ok := settings["MAX_HEADER_LIST_SIZE"]; ok {
		parts = append(parts, fmt.Sprintf("6:%d", val))
	}
	
	return fmt.Sprintf("%s|%d|0|m,a,p,s", strings.Join(parts, ";"), windowUpdate)
} 