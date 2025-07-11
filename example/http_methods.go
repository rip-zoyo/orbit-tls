package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/rip-zoyo/orbit-tls"
)

func main() {
	fmt.Println("🚀 Orbit TLS - HTTP Methods Demonstration")
	fmt.Println("=========================================")

	baseURL := "https://httpbin.org"
	client := orbit.Chrome138

	fmt.Println("\n📥 GET Request:")
	getResp, err := client.Get(baseURL + "/get")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", getResp.StatusCode)
	fmt.Printf("JA3 Hash: %s\n", getResp.GetJA3Hash())

	fmt.Println("\n📤 POST Request (JSON):")
	postData := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"active":   true,
		"score":    95.5,
	}

	postResp, err := client.PostJSON(baseURL+"/post", postData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", postResp.StatusCode)

	fmt.Println("\n📝 POST Request (Form Data):")
	formData := url.Values{
		"name":    {"John Doe"},
		"email":   {"john@example.com"},
		"message": {"Hello from Orbit TLS!"},
	}

	formResp, err := client.Post(baseURL+"/post", formData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", formResp.StatusCode)

	fmt.Println("\n🔧 PUT Request:")
	putData := map[string]interface{}{
		"id":     123,
		"name":   "Updated Item",
		"status": "active",
	}

	putResp, err := client.Put(baseURL+"/put", putData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", putResp.StatusCode)

	fmt.Println("\n🩹 PATCH Request:")
	patchData := map[string]interface{}{
		"status": "modified",
	}

	patchResp, err := client.Patch(baseURL+"/patch", patchData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", patchResp.StatusCode)

	fmt.Println("\n🗑️ DELETE Request:")
	deleteResp, err := client.Delete(baseURL + "/delete")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", deleteResp.StatusCode)

	fmt.Println("\n🔍 HEAD Request:")
	headResp, err := client.Head(baseURL + "/get")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", headResp.StatusCode)
	fmt.Printf("Content-Length: %s\n", headResp.Header.Get("Content-Length"))

	fmt.Println("\n⚙️ OPTIONS Request:")
	optionsResp, err := client.Options(baseURL + "/get")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %d\n", optionsResp.StatusCode)
	fmt.Printf("Allowed Methods: %s\n", optionsResp.Header.Get("Allow"))

	fmt.Println("\n🔐 Authenticated Requests:")
	authHeaders := map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"X-API-Key":     "secret-api-key-123",
	}

	authResp, err := client.Get(baseURL+"/bearer", authHeaders)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Auth Status: %d\n", authResp.StatusCode)

	fmt.Println("\n🌐 Request with Query Parameters:")
	params := map[string]string{
		"page":     "1",
		"limit":    "10",
		"sort":     "name",
		"filter":   "active",
	}

	options := &orbit.RequestOptions{
		Params: params,
	}

	paramResp, err := client.Request("GET", baseURL+"/get", nil, options)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Params Status: %d\n", paramResp.StatusCode)

	fmt.Println("\n🍪 Request with Cookies:")
	cookies := []*http.Cookie{
		{Name: "session_id", Value: "abc123def456"},
		{Name: "user_pref", Value: "dark_mode"},
	}

	cookieOptions := &orbit.RequestOptions{
		Cookies: cookies,
	}

	cookieResp, err := client.Request("GET", baseURL+"/cookies", nil, cookieOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Cookie Status: %d\n", cookieResp.StatusCode)

	fmt.Println("\n🎯 Multi-Browser Method Test:")
	browsers := map[string]*orbit.Client{
		"Chrome":  orbit.Chrome138,
		"Firefox": orbit.Firefox131,
		"Safari":  orbit.Safari18,
		"Brave":   orbit.Brave138,
	}

	testData := map[string]string{
		"browser": "test",
		"method":  "POST",
	}

	for name, browser := range browsers {
		resp, err := browser.PostJSON(baseURL+"/post", testData)
		if err != nil {
			log.Printf("❌ Error with %s: %v", name, err)
			continue
		}
		fmt.Printf("%s POST Status: %d\n", name, resp.StatusCode)
	}

	fmt.Println("\n✅ HTTP methods test completed!")
} 