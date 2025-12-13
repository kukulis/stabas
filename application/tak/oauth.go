package tak

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetOAuthToken obtains an OAuth access token using username/password
func GetOAuthToken(host, username, password, clientID string) (string, error) {
	url := fmt.Sprintf("https://%s:8446/oauth/token", host)

	// Prepare form data
	formData := fmt.Sprintf("grant_type=password&username=%s&password=%s&client_id=%s",
		username, password, clientID)

	// Configure HTTP client (skip TLS verification for testing)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // For testing only!
			},
		},
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(formData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OAuth failed: %s - %s", resp.Status, string(body))
	}

	// Parse JSON response manually (basic parsing)
	// For production, use encoding/json
	bodyStr := string(body)
	fmt.Printf("OAuth Response: %s\n", bodyStr)

	// Note: In production, properly unmarshal JSON response
	// This is a simplified example
	return "", fmt.Errorf("implement JSON parsing for production use")
}

// SendChatWithOAuth sends a chat message using OAuth token
func SendChatWithOAuth(host string, accessToken, cotXML string) error {
	url := fmt.Sprintf("https://%s:8446/Marti/api/cot", host)

	// Configure HTTP client (skip TLS verification for testing)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // For testing only!
			},
		},
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(cotXML))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/xml")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
