package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rip-zoyo/orbit-tls/fingerprint"
	"github.com/rip-zoyo/orbit-tls/profiles"
	"github.com/rip-zoyo/orbit-tls/tracking"
)

type Client struct {
	httpClient      *http.Client
	profile         *profiles.Profile
	headers         *OrderedHeaders
	tracker         *tracking.TLSTracker
	http2Tracker    *tracking.HTTP2Tracker
	lastFingerprint *fingerprint.Data
}

type Response struct {
	*http.Response
	Text        string           `json:"text"`
	Fingerprint *fingerprint.Data `json:"fingerprint,omitempty"`
}

type RequestOptions struct {
	Headers           map[string]string
	HeadersSlice      [][]string
	HeadersStringList []string
	HeadersJSON       string
	Cookies           []*http.Cookie
	Params            map[string]string
	Timeout           *int
}

type OrderedHeaders struct {
	headers []header
}

type header struct {
	name  string
	value string
}

var Chrome120 *Client
var Chrome131 *Client
var Chrome138 *Client

var Firefox121 *Client
var Firefox131 *Client

var Safari17 *Client
var Safari18 *Client
var SafariiOS *Client

var Edge120 *Client

var MullvadBrowser *Client
var ChromeAndroid *Client
var Opera115 *Client
var Brave131 *Client
var Brave138 *Client

func init() {
	var err error
	
	Chrome120, err = New("Chrome120")
	if err != nil {
		panic("Failed to create Chrome120 client: " + err.Error())
	}
	
	Chrome131, err = New("Chrome131")
	if err != nil {
		panic("Failed to create Chrome131 client: " + err.Error())
	}
	
	Chrome138, err = New("Chrome138")
	if err != nil {
		panic("Failed to create Chrome138 client: " + err.Error())
	}
	
	Firefox121, err = New("Firefox121")
	if err != nil {
		panic("Failed to create Firefox121 client: " + err.Error())
	}
	
	Firefox131, err = New("Firefox131")
	if err != nil {
		panic("Failed to create Firefox131 client: " + err.Error())
	}
	
	Safari17, err = New("Safari17")
	if err != nil {
		panic("Failed to create Safari17 client: " + err.Error())
	}
	
	Safari18, err = New("Safari18")
	if err != nil {
		panic("Failed to create Safari18 client: " + err.Error())
	}
	
	SafariiOS, err = New("SafariiOS")
	if err != nil {
		panic("Failed to create SafariiOS client: " + err.Error())
	}
	
	Edge120, err = New("Edge120")
	if err != nil {
		panic("Failed to create Edge120 client: " + err.Error())
	}
	
	MullvadBrowser, err = New("MullvadBrowser")
	if err != nil {
		panic("Failed to create MullvadBrowser client: " + err.Error())
	}
	
	ChromeAndroid, err = New("ChromeAndroid")
	if err != nil {
		panic("Failed to create ChromeAndroid client: " + err.Error())
	}
	
	Opera115, err = New("Opera115")
	if err != nil {
		panic("Failed to create Opera115 client: " + err.Error())
	}
	
	Brave131, err = New("Brave131")
	if err != nil {
		panic("Failed to create Brave131 client: " + err.Error())
	}
	
	Brave138, err = New("Brave138")
	if err != nil {
		panic("Failed to create Brave138 client: " + err.Error())
	}
}

func New(profileName string) (*Client, error) {
	profile, err := profiles.Get(profileName)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion:         profile.TLSVersion.Min,
		MaxVersion:         profile.TLSVersion.Max,
		InsecureSkipVerify: false,
		CipherSuites:       profile.CipherSuites,
		CurvePreferences:   profile.CurvePreferences,
		NextProtos:         profile.ALPNProtocols,
	}

	trackedDialer := tracking.NewTrackedDialer()
	
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return trackedDialer.DialTLS(network, addr, tlsConfig)
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	http2Tracker := tracking.NewHTTP2Tracker()
	
	defaultFrames := tracking.CreateHTTP2FrameForProfile(profile.Name)
	for _, frame := range defaultFrames {
		http2Tracker.TrackFrame(frame)
	}

	client := &Client{
		httpClient:   httpClient,
		profile:      profile,
		headers:      NewOrderedHeaders(),
		tracker:      tracking.GlobalTracker,
		http2Tracker: http2Tracker,
	}
	
	client.updateFingerprint()
	
	return client, nil
}

func (c *Client) Get(targetURL string, headers ...map[string]string) (*Response, error) {
	var opts *RequestOptions
	if len(headers) > 0 && headers[0] != nil {
		opts = &RequestOptions{Headers: headers[0]}
	}
	return c.Request("GET", targetURL, nil, opts)
}

func (c *Client) Post(targetURL string, body interface{}, headers ...map[string]string) (*Response, error) {
	var opts *RequestOptions
	if len(headers) > 0 && headers[0] != nil {
		opts = &RequestOptions{Headers: headers[0]}
	}
	return c.Request("POST", targetURL, body, opts)
}

func (c *Client) PostJSON(targetURL string, body interface{}, headers ...map[string]string) (*Response, error) {
	var jsonBody interface{}
	switch v := body.(type) {
	case string, []byte:
		jsonBody = v
	default:
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		jsonBody = jsonBytes
	}
	
	mergedHeaders := make(map[string]string)
	if len(headers) > 0 && headers[0] != nil {
		for k, v := range headers[0] {
			mergedHeaders[k] = v
		}
	}
	if _, exists := mergedHeaders["Content-Type"]; !exists {
		mergedHeaders["Content-Type"] = "application/json"
	}
	
	opts := &RequestOptions{Headers: mergedHeaders}
	return c.Request("POST", targetURL, jsonBody, opts)
}

func (c *Client) Put(targetURL string, body interface{}, headers ...map[string]string) (*Response, error) {
	var opts *RequestOptions
	if len(headers) > 0 && headers[0] != nil {
		opts = &RequestOptions{Headers: headers[0]}
	}
	return c.Request("PUT", targetURL, body, opts)
}

func (c *Client) Delete(targetURL string, headers ...map[string]string) (*Response, error) {
	var opts *RequestOptions
	if len(headers) > 0 && headers[0] != nil {
		opts = &RequestOptions{Headers: headers[0]}
	}
	return c.Request("DELETE", targetURL, nil, opts)
}

func (c *Client) Patch(targetURL string, body interface{}, headers ...map[string]string) (*Response, error) {
	var opts *RequestOptions
	if len(headers) > 0 && headers[0] != nil {
		opts = &RequestOptions{Headers: headers[0]}
	}
	return c.Request("PATCH", targetURL, body, opts)
}

func (c *Client) Head(targetURL string, headers ...map[string]string) (*Response, error) {
	var opts *RequestOptions
	if len(headers) > 0 && headers[0] != nil {
		opts = &RequestOptions{Headers: headers[0]}
	}
	return c.Request("HEAD", targetURL, nil, opts)
}

func (c *Client) Options(targetURL string, headers ...map[string]string) (*Response, error) {
	var opts *RequestOptions
	if len(headers) > 0 && headers[0] != nil {
		opts = &RequestOptions{Headers: headers[0]}
	}
	return c.Request("OPTIONS", targetURL, nil, opts)
}

func (c *Client) Request(method, targetURL string, body interface{}, options *RequestOptions) (*Response, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if options != nil && options.Params != nil {
		q := parsedURL.Query()
		for key, value := range options.Params {
			q.Add(key, value)
		}
		parsedURL.RawQuery = q.Encode()
	}

	var bodyReader io.Reader
	
	if body != nil {
		switch v := body.(type) {
		case string:
			bodyReader = strings.NewReader(v)
		case []byte:
			bodyReader = strings.NewReader(string(v))
		case io.Reader:
			bodyReader = v
		case url.Values:
			bodyReader = strings.NewReader(v.Encode())
		default:
			return nil, fmt.Errorf("unsupported body type: %T", body)
		}
	}

	req, err := http.NewRequest(method, parsedURL.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.applyAllHeaders(req, method, parsedURL, options)

	if options != nil && options.Cookies != nil {
		for _, cookie := range options.Cookies {
			req.AddCookie(cookie)
		}
	}

	if c.http2Tracker != nil {
		headersFrame := fingerprint.Frame{
			Type:     "HEADERS",
			StreamID: 1,
			Length:   uint32(tracking.EstimateHeadersSize(req)),
			Flags:    []string{"EndStream", "EndHeaders"},
			Headers:  tracking.ExtractHeadersList(req),
		}
		c.http2Tracker.TrackFrame(headersFrame)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	c.updateFingerprint()

	response := &Response{
		Response:    resp,
		Text:        string(responseBody),
		Fingerprint: c.lastFingerprint,
	}

	response.Body = io.NopCloser(strings.NewReader(response.Text))

	return response, nil
}

func (c *Client) updateFingerprint() {
	var details *tracking.ConnectionDetails
	connections := c.tracker.GetAllConnections()
	
	var latestTime time.Time
	for _, conn := range connections {
		if conn.ConnectedAt.After(latestTime) {
			latestTime = conn.ConnectedAt
			details = conn
		}
	}
	
	if details != nil {
		details.HTTP2Settings = c.http2Tracker.GetSettings()
		details.HTTP2Frames = c.http2Tracker.GetFrames()
		details.HTTP2WindowUpdate = c.http2Tracker.GetWindowSize()
		details.HTTP2Priority = c.http2Tracker.GetPriority()
	}
	
	c.lastFingerprint = tracking.GenerateFingerprintData(c.profile, details)
}

func (c *Client) applyAllHeaders(req *http.Request, method string, parsedURL *url.URL, opts *RequestOptions) {
	req.Header = make(http.Header)
	
	hasUserHeaders := c.hasUserHeaders(opts)
	
	if hasUserHeaders {
		for _, h := range c.headers.headers {
			req.Header.Set(h.name, h.value)
		}
		
		c.applyRequestHeaders(req, opts)
		
		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", c.profile.UserAgent)
		}
		if req.Header.Get("Accept") == "" {
			req.Header.Set("Accept", c.profile.Accept)
		}
		if req.Header.Get("Accept-Language") == "" {
			req.Header.Set("Accept-Language", c.profile.AcceptLanguage)
		}
		if req.Header.Get("Accept-Encoding") == "" {
			req.Header.Set("Accept-Encoding", c.profile.AcceptEncoding)
		}
	} else {
		profileHeaders := c.getProfileHeaders(method, parsedURL)
		for _, headerName := range c.profile.HeaderOrder {
			if strings.HasPrefix(headerName, ":") {
				continue
			}
			
			if value, exists := profileHeaders[headerName]; exists {
				req.Header.Set(headerName, value)
			}
		}
		
		for _, h := range c.headers.headers {
			req.Header.Set(h.name, h.value)
		}
	}
	
	req.Host = parsedURL.Host
}

func (c *Client) hasUserHeaders(opts *RequestOptions) bool {
	if opts == nil {
		return false
	}
	
	return (opts.Headers != nil && len(opts.Headers) > 0) ||
		   (opts.HeadersSlice != nil && len(opts.HeadersSlice) > 0) ||
		   (opts.HeadersStringList != nil && len(opts.HeadersStringList) > 0) ||
		   (opts.HeadersJSON != "" && opts.HeadersJSON != "{}")
}

func (c *Client) applyRequestHeaders(req *http.Request, opts *RequestOptions) {
	if opts == nil {
		return
	}
	
	if opts.Headers != nil {
		for key, value := range opts.Headers {
			req.Header.Set(key, value)
		}
	}
	
	if opts.HeadersSlice != nil {
		for _, pair := range opts.HeadersSlice {
			if len(pair) == 2 {
				req.Header.Set(pair[0], pair[1])
			}
		}
	}
	
	if opts.HeadersStringList != nil {
		for _, headerStr := range opts.HeadersStringList {
			parts := strings.SplitN(headerStr, ":", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				req.Header.Set(name, value)
			}
		}
	}
	
	if opts.HeadersJSON != "" && opts.HeadersJSON != "{}" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(opts.HeadersJSON), &headers); err == nil {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}
	}
}

func (c *Client) getProfileHeaders(method string, parsedURL *url.URL) map[string]string {
	headers := map[string]string{
		"User-Agent":      c.profile.UserAgent,
		"Accept":          c.profile.Accept,
		"Accept-Language": c.profile.AcceptLanguage,
		"Accept-Encoding": c.profile.AcceptEncoding,
	}

	switch c.profile.Name {
	case "Chrome120", "Edge120":
		headers["sec-ch-ua"] = `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`
		headers["sec-ch-ua-mobile"] = "?0"
		headers["sec-ch-ua-platform"] = `"Windows"`
		headers["DNT"] = "1"
		headers["Connection"] = "keep-alive"
		headers["Upgrade-Insecure-Requests"] = "1"
		
		switch method {
		case "GET":
			headers["Sec-Fetch-Site"] = "none"
			headers["Sec-Fetch-Mode"] = "navigate"
			headers["Sec-Fetch-User"] = "?1"
			headers["Sec-Fetch-Dest"] = "document"
		case "POST", "PUT", "PATCH":
			headers["Sec-Fetch-Site"] = "same-origin"
			headers["Sec-Fetch-Mode"] = "cors"
			headers["Sec-Fetch-Dest"] = "empty"
		default:
			headers["Sec-Fetch-Site"] = "same-origin"
			headers["Sec-Fetch-Mode"] = "cors"
			headers["Sec-Fetch-Dest"] = "empty"
		}
		
	case "Firefox121":
		headers["DNT"] = "1"
		headers["Connection"] = "keep-alive"
		headers["Upgrade-Insecure-Requests"] = "1"
		
	case "Safari17":
		headers["Connection"] = "keep-alive"
		headers["Upgrade-Insecure-Requests"] = "1"
	}

	if c.profile.SecHeaders != nil {
		for key, value := range c.profile.SecHeaders {
			headers[key] = value
		}
	}

	return headers
}

func (c *Client) GetJA3() string {
	return c.profile.JA3
}

func (c *Client) GetJA3Hash() string {
	return fingerprint.GenerateJA3Hash(c.profile.JA3)
}

func (c *Client) GetJA4() string {
	if c.lastFingerprint != nil {
		return c.lastFingerprint.JA4
	}
	
	tlsVersion, cipherSuites, extensions, supportedGroups, _, err := fingerprint.ParseJA3(c.profile.JA3)
	if err != nil {
		return ""
	}
	return fingerprint.GenerateJA4(tlsVersion, cipherSuites, extensions, supportedGroups, c.profile.ALPNProtocols)
}

func (c *Client) GetJA4R() string {
	if c.lastFingerprint != nil {
		return c.lastFingerprint.JA4_R
	}
	
	tlsVersion, cipherSuites, extensions, supportedGroups, _, err := fingerprint.ParseJA3(c.profile.JA3)
	if err != nil {
		return ""
	}
	return fingerprint.GenerateJA4R(tlsVersion, cipherSuites, extensions, supportedGroups, c.profile.ALPNProtocols)
}

func (c *Client) GetPeetPrint() string {
	if c.lastFingerprint != nil {
		return c.lastFingerprint.PeetPrint
	}
	return ""
}

func (c *Client) GetAkamaiFingerprint() string {
	if c.lastFingerprint != nil {
		return c.lastFingerprint.AkamaiFP
	}
	return ""
}

func (c *Client) GetClientRandom() string {
	if c.lastFingerprint != nil {
		return c.lastFingerprint.ClientRandom
	}
	return ""
}

func (c *Client) GetSessionID() string {
	if c.lastFingerprint != nil {
		return c.lastFingerprint.SessionID
	}
	return ""
}

func (c *Client) SetHeader(name, value string) {
	c.headers.Set(name, value)
}

func (c *Client) SetHeaders(headers map[string]string) {
	c.headers.SetMultiple(headers)
}

func (c *Client) SetHeadersFromSlice(headers [][]string) error {
	return c.headers.SetFromSlice(headers)
}

func (c *Client) SetHeadersFromStringSlice(headers []string) error {
	return c.headers.SetFromStringSlice(headers)
}

func (c *Client) SetHeadersFromJSON(jsonStr string) error {
	return c.headers.SetFromJSON(jsonStr)
}

func (c *Client) GetHeader(name string) string {
	return c.headers.Get(name)
}

func (c *Client) GetHeaders() map[string]string {
	return c.headers.GetAll()
}

func (c *Client) GetHeadersOrdered() [][]string {
	return c.headers.GetAllOrdered()
}

func (c *Client) DelHeader(name string) {
	c.headers.Del(name)
}

func (c *Client) ClearHeaders() {
	c.headers.Clear()
}

func NewOrderedHeaders() *OrderedHeaders {
	return &OrderedHeaders{
		headers: make([]header, 0),
	}
}

func (oh *OrderedHeaders) Set(name, value string) {
	name = strings.ToLower(name)
	
	for i, h := range oh.headers {
		if h.name == name {
			oh.headers[i].value = value
			return
		}
	}
	
	oh.headers = append(oh.headers, header{name: name, value: value})
}

func (oh *OrderedHeaders) Get(name string) string {
	name = strings.ToLower(name)
	for _, h := range oh.headers {
		if h.name == name {
			return h.value
		}
	}
	return ""
}

func (oh *OrderedHeaders) Del(name string) {
	name = strings.ToLower(name)
	for i, h := range oh.headers {
		if h.name == name {
			oh.headers = append(oh.headers[:i], oh.headers[i+1:]...)
			return
		}
	}
}

func (oh *OrderedHeaders) SetMultiple(headers map[string]string) {
	for name, value := range headers {
		oh.Set(name, value)
	}
}

func (oh *OrderedHeaders) SetFromSlice(headers [][]string) error {
	for _, pair := range headers {
		if len(pair) != 2 {
			return fmt.Errorf("invalid header pair: expected [name, value], got %v", pair)
		}
		oh.Set(pair[0], pair[1])
	}
	return nil
}

func (oh *OrderedHeaders) SetFromStringSlice(headers []string) error {
	for _, headerStr := range headers {
		parts := strings.SplitN(headerStr, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s (expected 'name: value')", headerStr)
		}
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		oh.Set(name, value)
	}
	return nil
}

func (oh *OrderedHeaders) SetFromJSON(jsonStr string) error {
	var headers map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &headers); err != nil {
		return fmt.Errorf("failed to parse JSON headers: %w", err)
	}
	oh.SetMultiple(headers)
	return nil
}

func (oh *OrderedHeaders) Clear() {
	oh.headers = make([]header, 0)
}

func (oh *OrderedHeaders) GetAll() map[string]string {
	result := make(map[string]string)
	for _, h := range oh.headers {
		result[h.name] = h.value
	}
	return result
}

func (oh *OrderedHeaders) GetAllOrdered() [][]string {
	result := make([][]string, len(oh.headers))
	for i, h := range oh.headers {
		result[i] = []string{h.name, h.value}
	}
	return result
}

func (r *Response) GetJA3() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.JA3
	}
	return ""
}

func (r *Response) GetJA3Hash() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.JA3Hash
	}
	return ""
}

func (r *Response) GetJA4() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.JA4
	}
	return ""
}

func (r *Response) GetJA4R() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.JA4_R
	}
	return ""
}

func (r *Response) GetPeetPrint() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.PeetPrint
	}
	return ""
}

func (r *Response) GetAkamaiFingerprint() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.AkamaiFP
	}
	return ""
}

func (r *Response) GetClientRandom() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.ClientRandom
	}
	return ""
}

func (r *Response) GetSessionID() string {
	if r.Fingerprint != nil {
		return r.Fingerprint.SessionID
	}
	return ""
} 