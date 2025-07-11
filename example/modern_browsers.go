package main

import (
	"fmt"
	"log"

	"github.com/rip-zoyo/orbit-tls"
)

func main() {
	fmt.Println("🚀 Orbit TLS - Modern Browsers with Advanced Features")
	fmt.Println("=====================================================")

	testURL := "https://tls.peet.ws/api/all"

	fmt.Println("\n🔮 Chrome 138 - Latest Features:")
	resp1, err := orbit.Chrome138.Get(testURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("JA3: %s\n", resp1.GetJA3())
	fmt.Printf("JA4: %s\n", resp1.GetJA4())
	fmt.Printf("JA4_R: %s\n", resp1.GetJA4R())
	fmt.Printf("PeetPrint: %s\n", resp1.GetPeetPrint())
	fmt.Printf("Akamai FP: %s\n", resp1.GetAkamaiFingerprint())
	fmt.Printf("Status: %d\n", resp1.StatusCode)

	fmt.Println("\n🦁 Brave 138 - Privacy-First with Post-Quantum:")
	resp2, err := orbit.Brave138.Get(testURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("JA3: %s\n", resp2.GetJA3())
	fmt.Printf("JA4: %s\n", resp2.GetJA4())
	fmt.Printf("JA4_R: %s\n", resp2.GetJA4R())
	fmt.Printf("PeetPrint: %s\n", resp2.GetPeetPrint())
	fmt.Printf("Akamai FP: %s\n", resp2.GetAkamaiFingerprint())
	fmt.Printf("Status: %d\n", resp2.StatusCode)

	fmt.Println("\n🔍 Feature Comparison:")
	fmt.Printf("Chrome138 vs Brave138 JA3 match: %t\n", resp1.GetJA3() == resp2.GetJA3())
	fmt.Printf("Chrome138 vs Brave138 JA4 match: %t\n", resp1.GetJA4() == resp2.GetJA4())

	fmt.Println("\n🧪 Testing Brave-specific privacy headers:")
	privacyHeaders := map[string]string{
		"DNT":            "1",
		"Sec-GPC":        "1",
		"X-Forwarded-For": "127.0.0.1",
	}

	resp3, err := orbit.Brave138.Get("https://httpbin.org/headers", privacyHeaders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Brave with privacy headers Status: %d\n", resp3.StatusCode)

	fmt.Println("\n📱 Mobile vs Desktop Chrome:")
	respDesktop, err := orbit.Chrome138.Get("https://httpbin.org/user-agent")
	if err != nil {
		log.Fatal(err)
	}

	respMobile, err := orbit.ChromeAndroid.Get("https://httpbin.org/user-agent")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Desktop Chrome JA3: %s\n", respDesktop.GetJA3Hash())
	fmt.Printf("Mobile Chrome JA3:  %s\n", respMobile.GetJA3Hash())
	fmt.Printf("Desktop vs Mobile match: %t\n", respDesktop.GetJA3Hash() == respMobile.GetJA3Hash())

	fmt.Println("\n🆚 Browser Generation Comparison:")
	browsers := map[string]*orbit.Client{
		"Chrome 120":  orbit.Chrome120,
		"Chrome 131":  orbit.Chrome131, 
		"Chrome 138":  orbit.Chrome138,
		"Brave 131":   orbit.Brave131,
		"Brave 138":   orbit.Brave138,
	}

	for name, client := range browsers {
		resp, err := client.Get("https://httpbin.org/headers")
		if err != nil {
			log.Printf("❌ Error with %s: %v", name, err)
			continue
		}
		fmt.Printf("%-12s JA3 Hash: %s\n", name, resp.GetJA3Hash())
	}

	fmt.Println("\n🔐 Advanced TLS Features Test:")
	fmt.Println("Testing post-quantum cryptography support...")
	
	advancedBrowsers := []*orbit.Client{orbit.Chrome138, orbit.Brave138}
	names := []string{"Chrome138", "Brave138"}

	for i, browser := range advancedBrowsers {
		resp, err := browser.Get("https://tls.peet.ws/api/all")
		if err != nil {
			log.Printf("❌ Error with %s: %v", names[i], err)
			continue
		}
		
		fmt.Printf("\n%s Advanced Features:\n", names[i])
		fmt.Printf("  Client Random: %s\n", resp.GetClientRandom())
		fmt.Printf("  Session ID: %s\n", resp.GetSessionID())
		fmt.Printf("  Full JA3: %s\n", resp.GetJA3())
	}

	fmt.Println("\n✅ Modern browsers test completed!")
} 