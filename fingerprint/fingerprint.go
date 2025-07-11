package fingerprint

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Data struct {
	TLSVersion           string        `json:"tls_version"`
	TLSVersionRecord     string        `json:"tls_version_record"`
	TLSVersionNegotiated string        `json:"tls_version_negotiated"`
	CipherSuites         []string      `json:"cipher_suites"`
	Extensions           []Extension   `json:"extensions"`
	SupportedGroups      []string      `json:"supported_groups"`
	SignatureAlgorithms  []string      `json:"signature_algorithms"`
	JA3                  string        `json:"ja3"`
	JA3Hash              string        `json:"ja3_hash"`
	JA4                  string        `json:"ja4"`
	JA4_R                string        `json:"ja4_r"`
	PeetPrint            string        `json:"peet_print"`
	PeetPrintHash        string        `json:"peet_print_hash"`
	AkamaiFP             string        `json:"akamai_fingerprint"`
	AkamaiFPHash         string        `json:"akamai_fingerprint_hash"`
	ClientRandom         string        `json:"client_random"`
	SessionID            string        `json:"session_id"`
	HTTP2                *HTTP2Data    `json:"http2,omitempty"`
}

type Extension struct {
	Type uint16 `json:"type"`
	Name string `json:"name"`
}

type HTTP2Data struct {
	Settings              map[string]uint32 `json:"settings"`
	WindowUpdate          uint32           `json:"window_update"`
	HeaderPriority        *HeaderPriority  `json:"header_priority,omitempty"`
	SentFrames            []Frame          `json:"sent_frames"`
	ConnectionPreface     string           `json:"connection_preface"`
	AkamaiFingerprint     string           `json:"akamai_fingerprint"`
	AkamaiFingerprintHash string           `json:"akamai_fingerprint_hash"`
}

type HeaderPriority struct {
	Dependency uint32 `json:"dependency"`
	Exclusive  bool   `json:"exclusive"`
	Weight     uint8  `json:"weight"`
}

type Frame struct {
	Type     string            `json:"type"`
	StreamID uint32            `json:"stream_id"`
	Length   uint32            `json:"length"`
	Flags    []string          `json:"flags"`
	Headers  []string          `json:"headers,omitempty"`
	Settings map[string]uint32 `json:"settings,omitempty"`
}

var cipherSuiteNames = map[uint16]string{
	0x1301: "TLS_AES_128_GCM_SHA256",
	0x1302: "TLS_AES_256_GCM_SHA384",
	0x1303: "TLS_CHACHA20_POLY1305_SHA256",
	0x1304: "TLS_AES_128_CCM_SHA256",
	0x1305: "TLS_AES_128_CCM_8_SHA256",
	0x0035: "TLS_RSA_WITH_AES_256_CBC_SHA",
	0x002f: "TLS_RSA_WITH_AES_128_CBC_SHA",
	0x000a: "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	0xc027: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	0xc028: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
	0xc013: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	0xc014: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	0xc02f: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	0xc030: "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	0xcca8: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
	0xcca9: "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
	0xc02b: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	0xc02c: "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	0xc023: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
	0xc024: "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384",
	0xc009: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	0xc00a: "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
}

var extensionNames = map[uint16]string{
	0:     "server_name",
	1:     "max_fragment_length",
	5:     "status_request",
	10:    "supported_groups",
	11:    "ec_point_formats",
	13:    "signature_algorithms",
	16:    "application_layer_protocol_negotiation",
	18:    "signed_certificate_timestamp",
	21:    "padding",
	23:    "session_ticket",
	27:    "compressed_certificate",
	35:    "session_ticket_tls",
	43:    "supported_versions",
	45:    "psk_key_exchange_modes",
	51:    "key_share",
	17513: "application_settings",
	65281: "renegotiation_info",
}

func GenerateJA3Hash(ja3 string) string {
	hash := md5.Sum([]byte(ja3))
	return hex.EncodeToString(hash[:])
}

func GenerateJA4(tlsVersion uint16, cipherSuites, extensions, supportedGroups []uint16, alpnProtocols []string) string {
	ver := "13"
	if tlsVersion == 0x0303 {
		ver = "12"
	} else if tlsVersion == 0x0302 {
		ver = "11"
	}

	proto := "d"
	if len(alpnProtocols) > 0 {
		if contains(alpnProtocols, "h2") {
			proto = "h2"
		} else if contains(alpnProtocols, "http/1.1") {
			proto = "h1"
		}
	}

	cipherHex := make([]string, len(cipherSuites))
	for i, c := range cipherSuites {
		cipherHex[i] = fmt.Sprintf("%04x", c)
	}
	sort.Strings(cipherHex)
	cipherStr := strings.Join(cipherHex, ",")

	extHex := make([]string, len(extensions))
	for i, e := range extensions {
		extHex[i] = fmt.Sprintf("%04x", e)
	}
	sort.Strings(extHex)
	extStr := strings.Join(extHex, ",")

	return fmt.Sprintf("t%s%s%s_%s_%s", ver, proto, "1516", cipherStr, extStr)
}

func GenerateJA4R(tlsVersion uint16, cipherSuites, extensions, supportedGroups []uint16, alpnProtocols []string) string {
	ver := "13"
	if tlsVersion == 0x0303 {
		ver = "12"
	} else if tlsVersion == 0x0302 {
		ver = "11"
	}

	proto := "d"
	if len(alpnProtocols) > 0 {
		if contains(alpnProtocols, "h2") {
			proto = "h2"
		} else if contains(alpnProtocols, "http/1.1") {
			proto = "h1"
		}
	}

	cipherHex := make([]string, len(cipherSuites))
	for i, c := range cipherSuites {
		cipherHex[i] = fmt.Sprintf("%04x", c)
	}
	cipherStr := strings.Join(cipherHex, ",")

	extHex := make([]string, len(extensions))
	for i, e := range extensions {
		extHex[i] = fmt.Sprintf("%04x", e)
	}
	extStr := strings.Join(extHex, ",")

	groupHex := make([]string, len(supportedGroups))
	for i, g := range supportedGroups {
		groupHex[i] = fmt.Sprintf("%04x", g)
	}
	groupStr := strings.Join(groupHex, ",")

	return fmt.Sprintf("t%s%s_%s_%s_%s", ver, proto, cipherStr, extStr, groupStr)
}

func GeneratePeetPrint(tlsVersionNegotiated, tlsVersionRecord uint16, cipherSuites, extensions, supportedGroups, signatureAlgorithms []uint16) string {
	parts := []string{
		fmt.Sprintf("%d-%d", tlsVersionNegotiated, tlsVersionRecord),
		formatUint16Slice(supportedGroups),
		formatUint16Slice(signatureAlgorithms),
		formatUint16Slice(extensions),
		formatUint16Slice(cipherSuites),
	}
	return strings.Join(parts, "|")
}

func GeneratePeetPrintHash(peetPrint string) string {
	hash := sha256.Sum256([]byte(peetPrint))
	return hex.EncodeToString(hash[:])
}

func ParseJA3(ja3 string) (uint16, []uint16, []uint16, []uint16, []uint16, error) {
	parts := strings.Split(ja3, ",")
	if len(parts) != 5 {
		return 0, nil, nil, nil, nil, fmt.Errorf("invalid JA3 format")
	}

	tlsVersion, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	cipherSuites, err := parseUint16List(parts[1])
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	extensions, err := parseUint16List(parts[2])
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	supportedGroups, err := parseUint16List(parts[3])
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	ecPointFormats, err := parseUint16List(parts[4])
	if err != nil {
		return 0, nil, nil, nil, nil, err
	}

	return uint16(tlsVersion), cipherSuites, extensions, supportedGroups, ecPointFormats, nil
}

func GetCipherSuiteName(id uint16) string {
	if name, exists := cipherSuiteNames[id]; exists {
		return name
	}
	return fmt.Sprintf("UNKNOWN_CIPHER_0x%04X", id)
}

func GetExtensionName(id uint16) string {
	if name, exists := extensionNames[id]; exists {
		return name
	}
	return fmt.Sprintf("UNKNOWN_EXTENSION_%d", id)
}

func parseUint16List(s string) ([]uint16, error) {
	if s == "" {
		return []uint16{}, nil
	}
	
	parts := strings.Split(s, "-")
	result := make([]uint16, len(parts))
	
	for i, part := range parts {
		val, err := strconv.ParseUint(part, 10, 16)
		if err != nil {
			return nil, err
		}
		result[i] = uint16(val)
	}
	
	return result, nil
}

func formatUint16Slice(slice []uint16) string {
	if len(slice) == 0 {
		return ""
	}
	
	strs := make([]string, len(slice))
	for i, v := range slice {
		strs[i] = strconv.FormatUint(uint64(v), 10)
	}
	
	return strings.Join(strs, "-")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
} 