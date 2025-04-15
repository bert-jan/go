package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// Target URL
	targetURL := fmt.Sprintf("https://httpbin.org/post")

	// Data params
	data := url.Values{}
	data.Set("us", "er")
	data.Set("pa", "ss")

	// Create POST request
	resp, err := http.Post(
		targetURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		fmt.Println("Error in request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:")
	fmt.Println(string(body))
}
