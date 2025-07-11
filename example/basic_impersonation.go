package main

import (
	"fmt"
	"log"

	"github.com/rip-zoyo/orbit-tls"
)

func main() {
	fmt.Println("🚀 Orbit TLS - Basic Browser Impersonation")
	fmt.Println("==========================================")

	testURL := "https://httpbin.org/user-agent"

	browsers := map[string]*orbit.Client{
		"Chrome 138":       orbit.Chrome138,
		"Firefox 131":     orbit.Firefox131,
		"Safari 18":       orbit.Safari18,
		"Edge 120":        orbit.Edge120,
		"Brave 138":       orbit.Brave138,
		"Chrome Android":  orbit.ChromeAndroid,
		"Safari iOS":      orbit.SafariiOS,
		"Mullvad Browser": orbit.MullvadBrowser,
	}

	for name, client := range browsers {
		fmt.Printf("\n🌐 Testing %s:\n", name)
		
		resp, err := client.Get(testURL)
		if err != nil {
			log.Printf("❌ Error with %s: %v", name, err)
			continue
		}

		fmt.Printf("   User-Agent detected: %s\n", extractUserAgent(resp.Text))
		fmt.Printf("   JA3 Hash: %s\n", resp.GetJA3Hash())
		fmt.Printf("   Status: %d\n", resp.StatusCode)
	}

	fmt.Println("\n✅ Browser impersonation test completed!")
}

func extractUserAgent(response string) string {
	start := "\"user-agent\": \""
	end := "\""
	
	startIdx := findInString(response, start)
	if startIdx == -1 {
		return "Not found"
	}
	
	startIdx += len(start)
	endIdx := findInStringFrom(response, end, startIdx)
	if endIdx == -1 {
		return "Not found"
	}
	
	return response[startIdx:endIdx]
}

func findInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func findInStringFrom(s, substr string, start int) int {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
} 