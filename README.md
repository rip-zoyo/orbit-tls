# Orbit TLS

A comprehensive TLS fingerprinting library for Go that accurately emulates browser behavior. Supports multiple fingerprinting methods including JA3, JA4, PeetPrint, and Akamai HTTP/2 fingerprinting.

## Installation

```bash
go get github.com/rip-zoyo/orbit-tls
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    orbit "github.com/rip-zoyo/orbit-tls"
)

func main() {
    client := orbit.Chrome138
    
    resp, err := client.Get("https://httpbin.org/user-agent")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Status: %d\n", resp.StatusCode)
    fmt.Printf("Response: %s\n", resp.Text)
    fmt.Printf("JA3: %s\n", resp.GetJA3())
}
```

## Browser Profiles

### Chrome Versions
| Profile | Version | JA3 Hash | User Agent |
|---------|---------|----------|------------|
| `Chrome120` | Chrome 120 | `cd08e31494f9531f560d64c695473da9` | Windows NT 10.0; Win64; x64 |
| `Chrome131` | Chrome 131 | `cd08e31494f9531f560d64c695473da9` | Windows NT 10.0; Win64; x64 |
| `Chrome138` | Chrome 138 | `773906b0efdcc3d6c1322c844389ae0e` | Windows NT 10.0; Win64; x64 |
| `ChromeAndroid` | Chrome 131 Mobile | `cd08e31494f9531f560d64c695473da9` | Android 14; SM-G998B |

### Firefox Versions
| Profile | Version | JA3 Hash | User Agent |
|---------|---------|----------|------------|
| `Firefox121` | Firefox 121 | `51c64c77e60f3980eea90869b68c58a8` | Windows NT 10.0; Win64; x64 |
| `Firefox131` | Firefox 131 | `51c64c77e60f3980eea90869b68c58a8` | Windows NT 10.0; Win64; x64 |

### Safari Versions
| Profile | Version | JA3 Hash | User Agent |
|---------|---------|----------|------------|
| `Safari17` | Safari 17.2.1 | `2a34309f6e0c6a6c4c9b9d79e7f8a7e5` | macOS 10.15.7 |
| `Safari18` | Safari 18.1.1 | `2a34309f6e0c6a6c4c9b9d79e7f8a7e5` | macOS 10.15.7 |
| `SafariiOS` | Safari iOS 17.1 | `2a34309f6e0c6a6c4c9b9d79e7f8a7e5` | iPhone; CPU iPhone OS 17_1 |
| `SafariiOS18` | Safari iOS 18.1 | `2a34309f6e0c6a6c4c9b9d79e7f8a7e5` | iPhone; CPU iPhone OS 18_1 |

### Other Browsers
| Profile | Version | JA3 Hash | User Agent |
|---------|---------|----------|------------|
| `Edge120` | Edge 120 | `cd08e31494f9531f560d64c695473da9` | Windows NT 10.0; Win64; x64 |
| `Brave131` | Brave 131 | `cd08e31494f9531f560d64c695473da9` | Windows NT 10.0; Win64; x64 |
| `Brave138` | Brave 138 | `773906b0efdcc3d6c1322c844389ae0e` | Windows NT 10.0; Win64; x64 |
| `Opera115` | Opera 115 | `cd08e31494f9531f560d64c695473da9` | Windows NT 10.0; Win64; x64 |
| `MullvadBrowser` | Mullvad Browser | `51c64c77e60f3980eea90869b68c58a8` | Windows NT 10.0 |

## Features

### TLS Fingerprinting
- **JA3 & JA3 Hash**: Classic TLS fingerprinting with MD5 hashing
- **JA4 & JA4_R**: Next-generation TLS fingerprinting standard
- **PeetPrint**: Advanced TLS fingerprinting with detailed analysis
- **Client Random & Session ID**: TLS handshake data capture

### HTTP/2 Fingerprinting
- **Akamai Fingerprinting**: HTTP/2 frame analysis for bot detection
- **Settings Fingerprinting**: HTTP/2 SETTINGS frame analysis
- **Frame Tracking**: Complete HTTP/2 frame sequence tracking
- **Window Updates**: HTTP/2 flow control fingerprinting

### Browser Emulation
- **Accurate User Agents**: Real browser user agent strings
- **Header Ordering**: Browser-specific HTTP header ordering
- **Security Headers**: Proper sec-ch-ua, sec-fetch-* headers
- **HTTP/2 Settings**: Browser-specific HTTP/2 configuration

### Advanced Features
- **Custom Headers**: Override default headers while maintaining fingerprint
- **Connection Pooling**: Efficient connection reuse
- **Certificate Validation**: Full certificate chain verification
- **Real-time Analysis**: Live TLS and HTTP/2 data capture

## Usage Examples

### Basic Browser Impersonation

```go
client := orbit.Chrome138
resp, err := client.Get("https://httpbin.org/headers")
```

### Custom Headers

```go
client := orbit.Firefox131

// Set individual headers
client.SetHeader("Authorization", "Bearer token123")
client.SetHeader("X-API-Key", "secret")

// Set multiple headers from map
client.SetHeaders(map[string]string{
    "X-Custom-Header": "value1",
    "X-Another-Header": "value2",
})

// Set headers from slice of [name, value] pairs
client.SetHeadersFromSlice([][]string{
    {"X-Request-ID", "12345"},
    {"X-Trace-ID", "67890"},
})

// Set headers from string slice (format: "name: value")
client.SetHeadersFromStringSlice([]string{
    "Content-Type: application/json",
    "Accept: application/json",
})

// Set headers from JSON string
client.SetHeadersFromJSON(`{
    "Authorization": "Bearer token123",
    "X-API-Key": "secret"
}`)

resp, err := client.Get("https://api.example.com/data")

// Or pass headers directly in request options
resp, err = client.Request("GET", "https://api.example.com/data", nil, &orbit.RequestOptions{
    Headers: map[string]string{
        "Authorization": "Bearer token123",
        "X-API-Key": "secret",
    },
})

// Use HeadersSlice in request options
resp, err = client.Request("GET", "https://api.example.com/data", nil, &orbit.RequestOptions{
    HeadersSlice: [][]string{
        {"Authorization", "Bearer token123"},
        {"X-API-Key", "secret"},
    },
})

// Use HeadersStringList in request options
resp, err = client.Request("GET", "https://api.example.com/data", nil, &orbit.RequestOptions{
    HeadersStringList: []string{
        "Authorization: Bearer token123",
        "X-API-Key: secret",
    },
})

// Use HeadersJSON in request options
resp, err = client.Request("GET", "https://api.example.com/data", nil, &orbit.RequestOptions{
    HeadersJSON: `{
        "Authorization": "Bearer token123",
        "X-API-Key": "secret"
    }`,
})
```

### Multiple HTTP Methods

```go
client := orbit.Safari18

// GET request
resp, err := client.Get("https://httpbin.org/get")

// POST with JSON
payload := map[string]interface{}{
    "username": "john",
    "password": "secret",
}
resp, err = client.PostJSON("https://httpbin.org/post", payload)

// Custom request
resp, err = client.Request("PATCH", "https://httpbin.org/patch", nil)
```

### Fingerprint Analysis

```go
client := orbit.Chrome138
resp, err := client.Get("https://tls.peet.ws/api/all")

if err == nil {
    fmt.Printf("JA3: %s\n", resp.GetJA3())
    fmt.Printf("JA3 Hash: %s\n", resp.GetJA3Hash())
    fmt.Printf("JA4: %s\n", resp.GetJA4())
    fmt.Printf("PeetPrint: %s\n", resp.GetPeetPrint())
    fmt.Printf("Akamai: %s\n", resp.GetAkamaiFingerprint())
}
```

## Header Management

### Setting Headers

Orbit TLS provides multiple ways to set headers:

```go
client := orbit.Chrome138

// Individual headers
client.SetHeader("Authorization", "Bearer token123")

// Multiple headers from map
client.SetHeaders(map[string]string{
    "X-Custom": "value1",
    "X-Another": "value2",
})

// Headers from slice of pairs
client.SetHeadersFromSlice([][]string{
    {"X-Request-ID", "12345"},
    {"X-Trace-ID", "67890"},
})

// Headers from string slice
client.SetHeadersFromStringSlice([]string{
    "Content-Type: application/json",
    "Accept: application/json",
})

// Headers from JSON
client.SetHeadersFromJSON(`{
    "Authorization": "Bearer token123",
    "X-API-Key": "secret"
}`)
```

### Getting Headers

```go
// Get single header
auth := client.GetHeader("Authorization")

// Get all headers as map
headers := client.GetHeaders()

// Get headers with preserved order
ordered := client.GetHeadersOrdered() // Returns [][]string
```

### Managing Headers

```go
// Delete a header
client.DelHeader("X-Custom")

// Clear all custom headers
client.ClearHeaders()
```

## Custom Headers Behavior

When you set custom headers, Orbit TLS automatically:
- Skips profile-specific header ordering (except pseudo-headers)
- Doesn't add sec-ch-ua headers automatically
- Only adds basic headers if you don't provide them:
  - User-Agent (from profile)
  - Accept (from profile)
  - Accept-Language (from profile)
  - Accept-Encoding (from profile)

```go
client := orbit.Chrome138

// This will use full profile behavior
resp1, _ := client.Get("https://example.com")

// This will use custom header behavior
client.SetHeader("X-Custom", "value")
resp2, _ := client.Get("https://example.com")
```

## Header Lists

Orbit TLS supports multiple formats for setting headers in bulk, making it easy to work with different data sources and APIs.

### Supported Header List Formats

#### 1. Map (Key-Value Pairs)
```go
client := orbit.Chrome138

// Set multiple headers from map
headers := map[string]string{
    "Authorization": "Bearer token123",
    "X-API-Key": "secret",
    "Content-Type": "application/json",
}
client.SetHeaders(headers)

// Or pass directly in request
resp, err := client.Request("GET", "https://api.example.com", nil, &orbit.RequestOptions{
    Headers: headers,
})
```

#### 2. Slice of [Name, Value] Pairs
```go
// Perfect for ordered headers or when working with arrays
headersList := [][]string{
    {"Authorization", "Bearer token123"},
    {"X-API-Key", "secret"},
    {"X-Request-ID", "12345"},
    {"Content-Type", "application/json"},
}

client.SetHeadersFromSlice(headersList)

// Or in request options
resp, err := client.Request("GET", "https://api.example.com", nil, &orbit.RequestOptions{
    HeadersSlice: headersList,
})
```

#### 3. String Slice (Colon-Separated)
```go
// Great for configuration files or command-line arguments
headersStrings := []string{
    "Authorization: Bearer token123",
    "X-API-Key: secret",
    "Content-Type: application/json",
    "X-Custom-Header: custom-value",
}

client.SetHeadersFromStringSlice(headersStrings)

// Or in request options
resp, err := client.Request("GET", "https://api.example.com", nil, &orbit.RequestOptions{
    HeadersStringList: headersStrings,
})
```

#### 4. JSON String
```go
// Perfect for API responses or configuration files
jsonHeaders := `{
    "Authorization": "Bearer token123",
    "X-API-Key": "secret",
    "Content-Type": "application/json",
    "X-Timestamp": "2024-01-01T00:00:00Z"
}`

client.SetHeadersFromJSON(jsonHeaders)

// Or in request options
resp, err := client.Request("GET", "https://api.example.com", nil, &orbit.RequestOptions{
    HeadersJSON: jsonHeaders,
})
```

### Combining Multiple Header Sources

```go
client := orbit.Firefox131

// Set base headers
client.SetHeaders(map[string]string{
    "User-Agent": "Custom-Agent/1.0",
    "Accept": "application/json",
})

// Add more headers from different sources
client.SetHeadersFromSlice([][]string{
    {"Authorization", "Bearer token123"},
    {"X-Request-ID", "req-456"},
})

// Headers from JSON config
configHeaders := `{"X-Environment": "production", "X-Version": "1.2.3"}`
client.SetHeadersFromJSON(configHeaders)

// All headers are combined
fmt.Printf("All headers: %+v\n", client.GetHeaders())
```

### Error Handling

```go
// Invalid slice format (missing value)
err := client.SetHeadersFromSlice([][]string{
    {"Authorization", "Bearer token123"},
    {"Invalid"}, // Missing value
})
if err != nil {
    log.Printf("Error: %v", err) // "invalid header pair: expected [name, value]"
}

// Invalid string format (missing colon)
err = client.SetHeadersFromStringSlice([]string{
    "Authorization: Bearer token123",
    "Invalid-No-Colon", // Missing colon
})
if err != nil {
    log.Printf("Error: %v", err) // "invalid header format"
}

// Invalid JSON
err = client.SetHeadersFromJSON(`{invalid json}`)
if err != nil {
    log.Printf("Error: %v", err) // "failed to parse JSON headers"
}
```

### Priority and Precedence

When using multiple header sources in the same request, the priority order is:

1. **RequestOptions.HeadersJSON** (highest priority)
2. **RequestOptions.HeadersStringList**
3. **RequestOptions.HeadersSlice**
4. **RequestOptions.Headers**
5. **Client persistent headers** (SetHeader, SetHeaders, etc.)
6. **Profile default headers** (lowest priority)

```go
client := orbit.Chrome138
client.SetHeader("X-Source", "client")

resp, err := client.Request("GET", "https://example.com", nil, &orbit.RequestOptions{
    Headers: map[string]string{"X-Source": "request-headers"},
    HeadersJSON: `{"X-Source": "request-json"}`,
})
// Final header value will be "request-json" (highest priority)
```

### Use Cases

**API Integration:**
```go
// Headers from API documentation
apiHeaders := []string{
    "Authorization: Bearer " + token,
    "X-API-Version: 2.1",
    "Content-Type: application/json",
}
client.SetHeadersFromStringSlice(apiHeaders)
```

**Configuration Files:**
```go
// Load headers from JSON config
configData, _ := os.ReadFile("headers.json")
client.SetHeadersFromJSON(string(configData))
```

**Dynamic Headers:**
```go
// Build headers programmatically
var headersList [][]string
for key, value := range dynamicHeaders {
    headersList = append(headersList, []string{key, value})
}
client.SetHeadersFromSlice(headersList)
```

## Testing Endpoints

Common endpoints for testing fingerprints:
- `https://tls.peet.ws/api/all` - Complete TLS analysis
- `https://ja3er.com/json` - JA3 fingerprint testing
- `https://httpbin.org/headers` - HTTP header analysis
- `https://httpbin.org/user-agent` - User agent testing

## Examples

The `example/` directory contains comprehensive examples:

1. **basic_impersonation.go** - Test all browser profiles
2. **custom_headers.go** - Custom header functionality
3. **modern_browsers.go** - Latest browser features
4. **http_methods.go** - All HTTP methods
5. **fingerprint_comparison.go** - Detailed fingerprint analysis
6. **header_lists.go** - Header list functionality and management

Run examples:
```bash
cd example
go run basic_impersonation.go
```

## Profile Details

Each profile includes:
- Accurate TLS cipher suites and extensions
- Browser-specific HTTP/2 settings
- Proper header ordering and security headers
- Real user agent strings
- Correct TLS version support
- Signature algorithms and supported groups

## Technical Implementation

- **TLS 1.2/1.3 Support**: Full modern TLS support
- **HTTP/2 Multiplexing**: Complete HTTP/2 implementation
- **Connection Reuse**: Efficient connection pooling
- **Memory Efficient**: Minimal memory allocations
- **Thread Safe**: Concurrent request support

## Requirements

- Go 1.21 or higher
- `golang.org/x/crypto`
- `golang.org/x/net`

## Contributing

Contributions welcome! Please submit pull requests for:
- New browser profiles
- Additional fingerprinting methods
- Bug fixes and improvements
- Documentation updates