package main

import (
	"fmt"
	"log"

	"github.com/rip-zoyo/orbit-tls"
)

func main() {
	fmt.Println("🚀 Orbit TLS - Custom Headers Example")
	fmt.Println("====================================")

	testURL := "https://httpbin.org/headers"

	fmt.Println("\n📋 Testing without custom headers (full profile behavior):")
	resp1, err := orbit.Chrome138.Get(testURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", resp1.StatusCode)
	fmt.Println("Response preview:", resp1.Text[:200]+"...")

	fmt.Println("\n🔧 Testing with custom headers (overrides profile):")
	customHeaders := map[string]string{
		"User-Agent":      "MyCustomBot/1.0",
		"Authorization":   "Bearer secret-token-123",
		"X-API-Key":       "my-api-key",
		"X-Custom-Header": "custom-value",
		"Accept":          "application/json",
	}

	resp2, err := orbit.Chrome138.Get(testURL, customHeaders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", resp2.StatusCode)
	fmt.Println("Response preview:", resp2.Text[:300]+"...")

	fmt.Println("\n🎯 Comparing fingerprints:")
	fmt.Printf("Default Chrome138 JA3: %s\n", resp1.GetJA3Hash())
	fmt.Printf("Custom headers JA3:    %s\n", resp2.GetJA3Hash())
	fmt.Printf("JA3 fingerprints match: %t\n", resp1.GetJA3Hash() == resp2.GetJA3Hash())

	fmt.Println("\n🧪 Testing API-style requests:")
	apiHeaders := map[string]string{
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"X-API-Version":   "v1",
		"X-Request-ID":    "req-12345",
	}

	data := map[string]interface{}{
		"username": "testuser",
		"action":   "login",
		"timestamp": "2024-01-01T00:00:00Z",
	}

	resp3, err := orbit.Firefox131.PostJSON("https://httpbin.org/post", data, apiHeaders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API POST Status: %d\n", resp3.StatusCode)

	fmt.Println("\n🔄 Testing header inheritance:")
	partialHeaders := map[string]string{
		"Authorization": "Bearer partial-override",
	}

	resp4, err := orbit.Safari18.Get(testURL, partialHeaders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Partial override Status: %d\n", resp4.StatusCode)
	fmt.Printf("Safari18 JA3: %s\n", resp4.GetJA3Hash())

	fmt.Println("\n✅ Custom headers test completed!")
} 