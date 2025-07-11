# Orbit TLS Examples

This directory contains comprehensive examples demonstrating the capabilities of the Orbit TLS library. Each example showcases different aspects of browser impersonation and TLS fingerprinting.

## Examples Overview

### 🌐 `basic_impersonation.go`
**Basic Browser Impersonation**

Demonstrates how to impersonate different browsers and compare their detected user agents. This example shows:
- Testing multiple browser profiles
- User-Agent detection and comparison
- JA3 fingerprint extraction
- Error handling for different browser clients

```bash
go run basic_impersonation.go
```

### 🔧 `custom_headers.go`
**Custom Headers and Profile Override**

Shows how custom headers override default profile behavior. This example demonstrates:
- Default profile behavior (full fingerprinting)
- Custom header override functionality
- API-style requests with authentication headers
- Partial header override (inheriting some profile headers)
- JA3 fingerprint consistency regardless of headers

```bash
go run custom_headers.go
```

### 🔮 `modern_browsers.go`
**Modern Browser Features**

Showcases the latest browser profiles with advanced features. This example covers:
- Chrome 138 and Brave 138 with post-quantum cryptography support
- Encrypted Client Hello (ECH) features
- Privacy-focused browser headers (Brave's Sec-GPC)
- Mobile vs Desktop Chrome comparison
- Browser generation evolution analysis
- Advanced TLS features testing

```bash
go run modern_browsers.go
```

### 🌐 `http_methods.go`
**HTTP Methods Demonstration**

Comprehensive demonstration of all supported HTTP methods. This example includes:
- GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS requests
- JSON POST requests
- Form data submission
- Authentication with custom headers
- Query parameters handling
- Cookie support
- Multi-browser method testing

```bash
go run http_methods.go
```

### 🔬 `fingerprint_comparison.go`
**Advanced Fingerprint Analysis**

Deep dive into fingerprint comparison across all browser profiles. This example provides:
- Complete fingerprint analysis (JA3, JA4, PeetPrint, Akamai)
- Browser family grouping and comparison
- Version evolution analysis
- Mobile vs Desktop fingerprint differences
- Privacy browser analysis
- Statistical uniqueness metrics

```bash
go run fingerprint_comparison.go
```

## Key Features Demonstrated

### Browser Profiles Covered
- **Chrome**: 120, 131, 138 (including Android)
- **Firefox**: 121, 131
- **Safari**: 17, 18 (including iOS)
- **Edge**: 120
- **Brave**: 131, 138
- **Opera**: 115
- **Mullvad Browser**: Privacy-focused
- **Safari iOS**: Mobile Safari

### Fingerprinting Methods
- **JA3**: TLS client fingerprinting
- **JA4**: Next-generation TLS fingerprinting
- **JA4_R**: JA4 with additional entropy
- **PeetPrint**: Advanced TLS fingerprinting
- **Akamai Fingerprinting**: Bot detection evasion

### Advanced Features
- **Post-Quantum Cryptography**: Chrome 138 & Brave 138
- **Encrypted Client Hello**: Modern TLS privacy
- **Custom Header Override**: Complete control over HTTP headers
- **Privacy Headers**: Brave's Sec-GPC and other privacy features
- **Mobile Fingerprinting**: Realistic mobile browser simulation

## Running Examples

Make sure you have the Orbit TLS library properly installed:

```bash
go mod tidy
```

Then run any example:

```bash
cd example
go run <example_file>.go
```

## Testing URLs

The examples use several testing endpoints:
- `https://httpbin.org/*` - HTTP testing service
- `https://tls.peet.ws/api/all` - TLS fingerprinting analysis

## Notes

- All examples handle errors gracefully
- Network connectivity required for external API calls
- Some examples may take time due to multiple requests
- Fingerprint results may vary depending on the target server's capabilities

## Custom Usage

These examples serve as templates for your own implementations. You can:
- Copy and modify code snippets
- Combine different approaches
- Add your own testing endpoints
- Implement custom fingerprint analysis logic

For more advanced usage, refer to the main library documentation. 