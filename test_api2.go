package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const defaultAPIC = "https://localhost"
const apiLoginEndpoint = "/api/aaaLogin.json"

type LoginResponse struct {
	Imdata []struct {
		AaaLogin struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"aaaLogin"`
	} `json:"imdata"`
}

type Endpoint struct {
	Dn   string `json:"dn"`
	Mac  string `json:"mac"`
	Ip   string `json:"ip"`
}

func loginToAPIC(apic string, username string, password string) (string, error) {
	loginData := map[string]interface{}{
		"aaaUser": map[string]interface{}{
			"attributes": map[string]interface{}{
				"name": username,
				"pwd":  password,
			},
		},
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apic+apiLoginEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return "", err
	}

	if len(loginResponse.Imdata) == 0 {
		return "", fmt.Errorf("failed to retrieve token")
	}

	return loginResponse.Imdata[0].AaaLogin.Attributes.Token, nil
}

func getEndpointList(apic, token, tenant, appProfile, epg string) ([]Endpoint, error) {
	epgDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, appProfile, epg)
	url := fmt.Sprintf("%s/api/node/class/fvCEp.json?query-target-filter=eq(fvCEp.epgDn,\"%s\")", apic, epgDn)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", "APIC-cookie="+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get endpoints, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var endpoints []Endpoint
	for _, item := range result["imdata"].([]interface{}) {
		epData := item.(map[string]interface{})["fvCEp"].(map[string]interface{})["attributes"].(map[string]interface{})
		endpoint := Endpoint{
			Dn:   epData["dn"].(string),
			Mac:  epData["mac"].(string),
			Ip:   epData["ip"].(string),
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

func deleteEndpoint(apic, token, dn string) error {
	url := fmt.Sprintf("%s/api/node/mo/%s.json", apic, dn)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", "APIC-cookie="+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to delete endpoint, status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	var apic, tenant, appProfile, epg, macAddress string
	var deleteFlag bool

	// Command line flags
	flag.StringVar(&apic, "H", defaultAPIC, "APIC IP address (default: localhost)")
	flag.StringVar(&tenant, "t", "", "Tenant name (required)")
	flag.StringVar(&appProfile, "a", "", "Application Profile name (required)")
	flag.StringVar(&epg, "e", "", "EPG name (required)")
	flag.StringVar(&macAddress, "m", "", "MAC Address to search for (required)")
	flag.BoolVar(&deleteFlag, "delete", false, "Delete the endpoint if found")
	flag.Parse()

	// Input validation
	if tenant == "" || appProfile == "" || epg == "" || macAddress == "" {
		fmt.Println("Error: Tenant, Application Profile, EPG, and MAC Address are required.")
		os.Exit(1)
	}

	// Login to APIC and get the token
	username := "admin"  // Hard-coded, change as needed
	password := "your_password" // Hard-coded, change as needed
	token, err := loginToAPIC(apic, username, password)
	if err != nil {
		fmt.Println("Login failed:", err)
		os.Exit(1)
	}

	// Get the list of endpoints for the given EPG
	endpoints, err := getEndpointList(apic, token, tenant, appProfile, epg)
	if err != nil {
		fmt.Println("Error getting endpoint list:", err)
		os.Exit(1)
	}

	// Check if the MAC address is found
	var foundEndpoint *Endpoint
	for _, ep := range endpoints {
		if strings.ToLower(ep.Mac) == strings.ToLower(macAddress) {
			foundEndpoint = &ep
			break
		}
	}

	if foundEndpoint != nil {
		fmt.Printf("Endpoint with MAC address %s found.\n", macAddress)
		if deleteFlag {
			// Delete the endpoint
			err := deleteEndpoint(apic, token, foundEndpoint.Dn)
			if err != nil {
				fmt.Println("Error deleting endpoint:", err)
				os.Exit(1)
			}
			fmt.Printf("Endpoint with MAC address %s successfully deleted.\n", macAddress)
		}
	} else {
		fmt.Printf("Endpoint with MAC address %s not found.\n", macAddress)
	}
}
