package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const GithubClientID = "178c6fc778ccc68e1d6a" // GitHub CLI official client ID

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

// PerformGithubOAuth starts the device authorization flow and polls until the user approves.
func PerformGithubOAuth() (string, error) {
	// 1. Request device code
	reqBody := []byte(fmt.Sprintf("client_id=%s&scope=repo read:org", GithubClientID))
	req, err := http.NewRequest("POST", "https://github.com/login/device/code", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var deviceRes DeviceCodeResponse
	if err := json.Unmarshal(body, &deviceRes); err != nil {
		return "", fmt.Errorf("failed to parse GitHub response: %v", err)
	}

	if deviceRes.UserCode == "" {
		return "", fmt.Errorf("invalid response from GitHub")
	}

	// 2. Prompt user
	fmt.Printf("\n🔒 GitHub Authentication Required\n")
	fmt.Printf("1. Please open: %s\n", deviceRes.VerificationURI)
	fmt.Printf("2. Enter the code: %s\n", deviceRes.UserCode)
	fmt.Printf("\nWaiting for you to authorize (polling)... ")

	// 3. Poll for access token
	pollInterval := time.Duration(deviceRes.Interval) * time.Second
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	tokenReqBody := []byte(fmt.Sprintf("client_id=%s&device_code=%s&grant_type=urn:ietf:params:oauth:grant-type:device_code", GithubClientID, deviceRes.DeviceCode))

	for i := 0; i < 60; i++ { // Timeout after 5 minutes (60 * 5s)
		time.Sleep(pollInterval)

		tokenReq, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(tokenReqBody))
		tokenReq.Header.Set("Accept", "application/json")
		tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tokenResp, err := client.Do(tokenReq)
		if err != nil {
			continue
		}

		tokenBody, _ := io.ReadAll(tokenResp.Body)
		tokenResp.Body.Close()

		var accessRes AccessTokenResponse
		_ = json.Unmarshal(tokenBody, &accessRes)

		if accessRes.AccessToken != "" {
			fmt.Printf("✅ Success!\n")
			return accessRes.AccessToken, nil
		}

		if accessRes.Error != "authorization_pending" {
			return "", fmt.Errorf("GitHub authorization failed: %s", accessRes.Error)
		}
	}

	return "", fmt.Errorf("authentication timed out")
}
