package main

import (
	"fmt"
	"log"

	"github.com/rip-zoyo/orbit-tls"
)

func main() {
	fmt.Println("🚀 Orbit TLS - Fingerprint Comparison")
	fmt.Println("=====================================")

	testURL := "https://tls.peet.ws/api/all"

	browsers := map[string]*orbit.Client{
		"Chrome 120":       orbit.Chrome120,
		"Chrome 131":       orbit.Chrome131,
		"Chrome 138":       orbit.Chrome138,
		"Firefox 121":      orbit.Firefox121,
		"Firefox 131":      orbit.Firefox131,
		"Safari 17":        orbit.Safari17,
		"Safari 18":        orbit.Safari18,
		"Safari iOS":       orbit.SafariiOS,
		"Edge 120":         orbit.Edge120,
		"Brave 131":        orbit.Brave131,
		"Brave 138":        orbit.Brave138,
		"Chrome Android":   orbit.ChromeAndroid,
		"Opera 115":        orbit.Opera115,
		"Mullvad Browser":  orbit.MullvadBrowser,
	}

	fmt.Println("\n📊 Complete Fingerprint Analysis:")
	fmt.Println("═══════════════════════════════════")

	results := make(map[string]map[string]string)

	for name, client := range browsers {
		fmt.Printf("\n🔍 Testing %s...\n", name)
		
		resp, err := client.Get(testURL)
		if err != nil {
			log.Printf("❌ Error with %s: %v", name, err)
			continue
		}

		results[name] = map[string]string{
			"JA3":      resp.GetJA3(),
			"JA3Hash":  resp.GetJA3Hash(),
			"JA4":      resp.GetJA4(),
			"JA4_R":    resp.GetJA4R(),
			"PeetPrint": resp.GetPeetPrint(),
			"AkamaiFP": resp.GetAkamaiFingerprint(),
		}

		fmt.Printf("  Status: %d\n", resp.StatusCode)
		fmt.Printf("  JA3 Hash: %s\n", results[name]["JA3Hash"])
		fmt.Printf("  JA4: %s\n", results[name]["JA4"])
	}

	fmt.Println("\n🔬 Fingerprint Analysis:")
	fmt.Println("═══════════════════════")

	fmt.Println("\n🏷️ JA3 Hash Comparison:")
	ja3Groups := make(map[string][]string)
	for name, data := range results {
		hash := data["JA3Hash"]
		ja3Groups[hash] = append(ja3Groups[hash], name)
	}

	for hash, browsers := range ja3Groups {
		fmt.Printf("JA3 Hash %s:\n", hash[:16]+"...")
		for _, browser := range browsers {
			fmt.Printf("  - %s\n", browser)
		}
		fmt.Println()
	}

	fmt.Println("\n🎯 JA4 Comparison:")
	ja4Groups := make(map[string][]string)
	for name, data := range results {
		ja4 := data["JA4"]
		if ja4 != "" {
			ja4Groups[ja4] = append(ja4Groups[ja4], name)
		}
	}

	for ja4, browsers := range ja4Groups {
		fmt.Printf("JA4 %s:\n", ja4)
		for _, browser := range browsers {
			fmt.Printf("  - %s\n", browser)
		}
		fmt.Println()
	}

	fmt.Println("\n🔬 Browser Family Analysis:")
	chromiumBrowsers := []string{"Chrome 120", "Chrome 131", "Chrome 138", "Edge 120", "Brave 131", "Brave 138", "Chrome Android", "Opera 115"}
	firefoxBrowsers := []string{"Firefox 121", "Firefox 131", "Mullvad Browser"}
	safariBrowsers := []string{"Safari 17", "Safari 18", "Safari iOS"}

	fmt.Println("\n🏢 Chromium Family:")
	for _, browser := range chromiumBrowsers {
		if data, exists := results[browser]; exists {
			fmt.Printf("  %-15s JA3: %s\n", browser, data["JA3Hash"][:16]+"...")
		}
	}

	fmt.Println("\n🦊 Firefox Family:")
	for _, browser := range firefoxBrowsers {
		if data, exists := results[browser]; exists {
			fmt.Printf("  %-15s JA3: %s\n", browser, data["JA3Hash"][:16]+"...")
		}
	}

	fmt.Println("\n🍎 Safari Family:")
	for _, browser := range safariBrowsers {
		if data, exists := results[browser]; exists {
			fmt.Printf("  %-15s JA3: %s\n", browser, data["JA3Hash"][:16]+"...")
		}
	}

	fmt.Println("\n🚀 Modern vs Legacy Features:")
	fmt.Println("\n🔮 Advanced Features (Chrome 138 & Brave 138):")
	
	advancedBrowsers := []string{"Chrome 138", "Brave 138"}
	for _, browser := range advancedBrowsers {
		if data, exists := results[browser]; exists {
			fmt.Printf("\n%s:\n", browser)
			fmt.Printf("  Full JA3: %s\n", data["JA3"])
			fmt.Printf("  JA4_R:    %s\n", data["JA4_R"])
			fmt.Printf("  PeetPrint: %s\n", data["PeetPrint"])
			fmt.Printf("  Akamai FP: %s\n", data["AkamaiFP"])
		}
	}

	fmt.Println("\n📱 Mobile vs Desktop Comparison:")
	mobileDesktopPairs := map[string][2]string{
		"Chrome": {"Chrome 138", "Chrome Android"},
		"Safari": {"Safari 18", "Safari iOS"},
	}

	for family, pair := range mobileDesktopPairs {
		desktop, mobile := pair[0], pair[1]
		if desktopData, exists1 := results[desktop]; exists1 {
			if mobileData, exists2 := results[mobile]; exists2 {
				fmt.Printf("\n%s Family:\n", family)
				fmt.Printf("  Desktop (%s): %s\n", desktop, desktopData["JA3Hash"][:16]+"...")
				fmt.Printf("  Mobile  (%s): %s\n", mobile, mobileData["JA3Hash"][:16]+"...")
				match := desktopData["JA3Hash"] == mobileData["JA3Hash"]
				fmt.Printf("  Fingerprints match: %t\n", match)
			}
		}
	}

	fmt.Println("\n🔄 Version Evolution Analysis:")
	chromeVersions := []string{"Chrome 120", "Chrome 131", "Chrome 138"}
	firefoxVersions := []string{"Firefox 121", "Firefox 131"}
	braveVersions := []string{"Brave 131", "Brave 138"}

	families := map[string][]string{
		"Chrome":  chromeVersions,
		"Firefox": firefoxVersions,
		"Brave":   braveVersions,
	}

	for family, versions := range families {
		fmt.Printf("\n%s Evolution:\n", family)
		for _, version := range versions {
			if data, exists := results[version]; exists {
				fmt.Printf("  %-12s %s\n", version, data["JA3Hash"][:16]+"...")
			}
		}
	}

	fmt.Println("\n🎭 Privacy Browser Analysis:")
	privacyBrowsers := []string{"Brave 131", "Brave 138", "Mullvad Browser"}
	fmt.Println("\nPrivacy-focused browsers:")
	for _, browser := range privacyBrowsers {
		if data, exists := results[browser]; exists {
			fmt.Printf("  %-15s JA3: %s\n", browser, data["JA3Hash"][:16]+"...")
		}
	}

	fmt.Println("\n📈 Summary Statistics:")
	uniqueJA3 := len(ja3Groups)
	uniqueJA4 := len(ja4Groups)
	totalBrowsers := len(results)

	fmt.Printf("Total browsers tested: %d\n", totalBrowsers)
	fmt.Printf("Unique JA3 fingerprints: %d\n", uniqueJA3)
	fmt.Printf("Unique JA4 fingerprints: %d\n", uniqueJA4)
	fmt.Printf("JA3 uniqueness ratio: %.2f%%\n", float64(uniqueJA3)/float64(totalBrowsers)*100)
	fmt.Printf("JA4 uniqueness ratio: %.2f%%\n", float64(uniqueJA4)/float64(totalBrowsers)*100)

	fmt.Println("\n✅ Fingerprint comparison completed!")
} 