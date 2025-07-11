package main

import (
	"fmt"
	"log"
	orbit "github.com/rip-zoyo/orbit-tls"
)

func main() {
	fmt.Println("=== 🎯 Header Lists Example ===")
	
	client := orbit.Chrome138
	
	// Example 1: Set headers using map
	fmt.Println("\n1. Setting headers from map:")
	client.SetHeaders(map[string]string{
		"X-API-Key":    "secret123",
		"X-Request-ID": "req-12345",
		"X-Source":     "golang-client",
	})
	
	// Show what headers are set
	fmt.Printf("Current headers: %+v\n", client.GetHeaders())
	
	// Example 2: Set headers from slice of pairs
	fmt.Println("\n2. Setting headers from slice:")
	client.ClearHeaders()
	err := client.SetHeadersFromSlice([][]string{
		{"Authorization", "Bearer token123"},
		{"Content-Type", "application/json"},
		{"X-Custom-Header", "custom-value"},
	})
	if err != nil {
		log.Printf("Error setting headers from slice: %v", err)
	}
	
	fmt.Printf("Headers from slice: %+v\n", client.GetHeaders())
	
	// Example 3: Set headers from string slice
	fmt.Println("\n3. Setting headers from string slice:")
	client.ClearHeaders()
	err = client.SetHeadersFromStringSlice([]string{
		"X-App-Version: 1.2.3",
		"X-Platform: linux",
		"X-Environment: production",
	})
	if err != nil {
		log.Printf("Error setting headers from string slice: %v", err)
	}
	
	fmt.Printf("Headers from string slice: %+v\n", client.GetHeaders())
	
	// Example 4: Set headers from JSON
	fmt.Println("\n4. Setting headers from JSON:")
	client.ClearHeaders()
	jsonHeaders := `{
		"Authorization": "Bearer jwt-token",
		"X-Client-ID": "client-456",
		"X-Timestamp": "2024-01-01T00:00:00Z"
	}`
	err = client.SetHeadersFromJSON(jsonHeaders)
	if err != nil {
		log.Printf("Error setting headers from JSON: %v", err)
	}
	
	fmt.Printf("Headers from JSON: %+v\n", client.GetHeaders())
	
	// Example 5: Using headers in request options
	fmt.Println("\n5. Using headers in request options:")
	
	// Using map in request options
	resp, err := client.Request("GET", "https://httpbin.org/headers", nil, &orbit.RequestOptions{
		Headers: map[string]string{
			"X-Test-Header": "from-request-options",
			"X-Method":      "map",
		},
	})
	if err != nil {
		log.Printf("Error with map headers: %v", err)
	} else {
		fmt.Printf("Response with map headers: %d\n", resp.StatusCode)
	}
	
	// Using slice in request options
	resp, err = client.Request("GET", "https://httpbin.org/headers", nil, &orbit.RequestOptions{
		HeadersSlice: [][]string{
			{"X-Test-Header", "from-slice"},
			{"X-Method", "slice"},
		},
	})
	if err != nil {
		log.Printf("Error with slice headers: %v", err)
	} else {
		fmt.Printf("Response with slice headers: %d\n", resp.StatusCode)
	}
	
	// Using string list in request options
	resp, err = client.Request("GET", "https://httpbin.org/headers", nil, &orbit.RequestOptions{
		HeadersStringList: []string{
			"X-Test-Header: from-string-list",
			"X-Method: string-list",
		},
	})
	if err != nil {
		log.Printf("Error with string list headers: %v", err)
	} else {
		fmt.Printf("Response with string list headers: %d\n", resp.StatusCode)
	}
	
	// Using JSON in request options
	resp, err = client.Request("GET", "https://httpbin.org/headers", nil, &orbit.RequestOptions{
		HeadersJSON: `{
			"X-Test-Header": "from-json",
			"X-Method": "json"
		}`,
	})
	if err != nil {
		log.Printf("Error with JSON headers: %v", err)
	} else {
		fmt.Printf("Response with JSON headers: %d\n", resp.StatusCode)
	}
	
	// Example 6: Header management
	fmt.Println("\n6. Header management:")
	
	// Set some headers
	client.SetHeaders(map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
		"X-Header-3": "value3",
	})
	
	fmt.Printf("All headers: %+v\n", client.GetHeaders())
	fmt.Printf("Ordered headers: %+v\n", client.GetHeadersOrdered())
	
	// Get specific header
	header1 := client.GetHeader("X-Header-1")
	fmt.Printf("X-Header-1 value: %s\n", header1)
	
	// Delete a header
	client.DelHeader("X-Header-2")
	fmt.Printf("After deleting X-Header-2: %+v\n", client.GetHeaders())
	
	// Clear all headers
	client.ClearHeaders()
	fmt.Printf("After clearing all headers: %+v\n", client.GetHeaders())
	
	// Example 7: Error handling
	fmt.Println("\n7. Error handling:")
	
	// Invalid slice format
	err = client.SetHeadersFromSlice([][]string{
		{"incomplete"},
		{"X-Good-Header", "good-value"},
	})
	if err != nil {
		fmt.Printf("Expected error with invalid slice: %v\n", err)
	}
	
	// Invalid string format
	err = client.SetHeadersFromStringSlice([]string{
		"valid-header: valid-value",
		"invalid-header-no-colon",
	})
	if err != nil {
		fmt.Printf("Expected error with invalid string: %v\n", err)
	}
	
	// Invalid JSON
	err = client.SetHeadersFromJSON(`{invalid json}`)
	if err != nil {
		fmt.Printf("Expected error with invalid JSON: %v\n", err)
	}
	
	fmt.Println("\n🎉 Header lists example completed!")
} 