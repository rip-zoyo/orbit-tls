package tracking

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/http2/hpack"
	"github.com/rip-zoyo/orbit-tls/fingerprint"
)

type HTTP2Tracker struct {
	mu         sync.RWMutex
	frames     []fingerprint.Frame
	settings   map[string]uint32
	windowSize uint32
	priority   *fingerprint.HeaderPriority
}

func NewHTTP2Tracker() *HTTP2Tracker {
	return &HTTP2Tracker{
		frames:   make([]fingerprint.Frame, 0),
		settings: make(map[string]uint32),
	}
}

func (t *HTTP2Tracker) TrackFrame(frame fingerprint.Frame) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	t.frames = append(t.frames, frame)
	
	if frame.Type == "SETTINGS" && frame.Settings != nil {
		for key, value := range frame.Settings {
			t.settings[key] = value
		}
	}
	
	if frame.Type == "WINDOW_UPDATE" && frame.Length > 0 {
		t.windowSize = frame.Length
	}
	
	if frame.Type == "HEADERS" && frame.Settings != nil {
		if priority, ok := extractPriority(frame); ok {
			t.priority = priority
		}
	}
}

func (t *HTTP2Tracker) GetFrames() []fingerprint.Frame {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	result := make([]fingerprint.Frame, len(t.frames))
	copy(result, t.frames)
	return result
}

func (t *HTTP2Tracker) GetSettings() map[string]uint32 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	result := make(map[string]uint32)
	for k, v := range t.settings {
		result[k] = v
	}
	return result
}

func (t *HTTP2Tracker) GetWindowSize() uint32 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.windowSize
}

func (t *HTTP2Tracker) GetPriority() *fingerprint.HeaderPriority {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.priority
}

func CreateDefaultChromeFrames() []fingerprint.Frame {
	return []fingerprint.Frame{
		{
			Type:   "SETTINGS",
			Length: 24,
			Settings: map[string]uint32{
				"HEADER_TABLE_SIZE":      65536,
				"ENABLE_PUSH":            0,
				"MAX_CONCURRENT_STREAMS": 1000,
				"INITIAL_WINDOW_SIZE":    6291456,
				"MAX_FRAME_SIZE":         16384,
				"MAX_HEADER_LIST_SIZE":   262144,
			},
		},
		{
			Type:   "WINDOW_UPDATE",
			Length: 1073741824,
		},
	}
}

func CreateDefaultFirefoxFrames() []fingerprint.Frame {
	return []fingerprint.Frame{
		{
			Type:   "SETTINGS",
			Length: 24,
			Settings: map[string]uint32{
				"HEADER_TABLE_SIZE":      65536,
				"ENABLE_PUSH":            1,
				"MAX_CONCURRENT_STREAMS": 1000,
				"INITIAL_WINDOW_SIZE":    131072,
				"MAX_FRAME_SIZE":         16384,
				"MAX_HEADER_LIST_SIZE":   262144,
			},
		},
		{
			Type:   "WINDOW_UPDATE",
			Length: 12517376,
		},
	}
}

func CreateDefaultSafariFrames() []fingerprint.Frame {
	return []fingerprint.Frame{
		{
			Type:   "SETTINGS",
			Length: 20,
			Settings: map[string]uint32{
				"HEADER_TABLE_SIZE":      4096,
				"ENABLE_PUSH":            1,
				"MAX_CONCURRENT_STREAMS": 100,
				"INITIAL_WINDOW_SIZE":    2097152,
				"MAX_FRAME_SIZE":         16384,
				"MAX_HEADER_LIST_SIZE":   8192,
			},
		},
		{
			Type:   "WINDOW_UPDATE",
			Length: 2013265920,
		},
	}
}

func CreateHTTP2FrameForProfile(profileName string) []fingerprint.Frame {
	switch profileName {
	case "Chrome120", "Edge120":
		return CreateDefaultChromeFrames()
	case "Firefox121":
		return CreateDefaultFirefoxFrames()
	case "Safari17":
		return CreateDefaultSafariFrames()
	default:
		return CreateDefaultChromeFrames()
	}
}

func EstimateHeadersSize(req *http.Request) int {
	size := 0
	size += len(req.Method) + len(req.URL.Path) + len(req.Proto)
	
	for name, values := range req.Header {
		for _, value := range values {
			size += len(name) + len(value) + 4
		}
	}
	
	return size
}

func ExtractHeadersList(req *http.Request) []string {
	var headers []string
	
	headers = append(headers, fmt.Sprintf(":method: %s", req.Method))
	headers = append(headers, fmt.Sprintf(":path: %s", req.URL.Path))
	headers = append(headers, fmt.Sprintf(":scheme: %s", req.URL.Scheme))
	headers = append(headers, fmt.Sprintf(":authority: %s", req.Host))
	
	for name, values := range req.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", strings.ToLower(name), value))
		}
	}
	
	return headers
}

func GenerateAkamaiH2Fingerprint(settings map[string]uint32, windowUpdate uint32, priority *fingerprint.HeaderPriority) string {
	var settingsParts []string
	
	if val, ok := settings["HEADER_TABLE_SIZE"]; ok {
		settingsParts = append(settingsParts, fmt.Sprintf("1:%d", val))
	}
	if val, ok := settings["ENABLE_PUSH"]; ok {
		settingsParts = append(settingsParts, fmt.Sprintf("2:%d", val))
	}
	if val, ok := settings["MAX_CONCURRENT_STREAMS"]; ok {
		settingsParts = append(settingsParts, fmt.Sprintf("3:%d", val))
	}
	if val, ok := settings["INITIAL_WINDOW_SIZE"]; ok {
		settingsParts = append(settingsParts, fmt.Sprintf("4:%d", val))
	}
	if val, ok := settings["MAX_FRAME_SIZE"]; ok {
		settingsParts = append(settingsParts, fmt.Sprintf("5:%d", val))
	}
	if val, ok := settings["MAX_HEADER_LIST_SIZE"]; ok {
		settingsParts = append(settingsParts, fmt.Sprintf("6:%d", val))
	}
	
	priorityStr := "0"
	if priority != nil {
		if priority.Exclusive {
			priorityStr = "1"
		}
	}
	
	return fmt.Sprintf("%s|%d|%s|m,a,p,s", strings.Join(settingsParts, ";"), windowUpdate, priorityStr)
}

func DecodeHPACKHeaders(headerBlock []byte) []string {
	var headers []string
	decoder := hpack.NewDecoder(4096, nil)
	
	headerFields, err := decoder.DecodeFull(headerBlock)
	if err != nil {
		return []string{"user-agent: Go-http-client/2.0"}
	}
	
	for _, field := range headerFields {
		headers = append(headers, fmt.Sprintf("%s: %s", field.Name, field.Value))
	}
	
	return headers
}

func extractPriority(frame fingerprint.Frame) (*fingerprint.HeaderPriority, bool) {
	if frame.Type != "HEADERS" {
		return nil, false
	}
	
	for _, flag := range frame.Flags {
		if flag == "Priority" {
			return &fingerprint.HeaderPriority{
				Dependency: 0,
				Exclusive:  false,
				Weight:     16,
			}, true
		}
	}
	
	return nil, false
} 