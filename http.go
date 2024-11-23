package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendPayload(rule Rule, payload interface{}) error {
	if rule.Http == nil {
		return fmt.Errorf("no HTTP configuration provided in rule")
	}

	// Serialize the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize payload: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(rule.Http.Method, rule.Http.Url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	for _, header := range rule.Http.Headers {
		req.Header.Set(header.Name, header.Value)
	}

	// Set authentication
	if rule.Http.Auth != nil {
		switch rule.Http.Auth.Type {
		case "bearer":
			req.Header.Set("Authorization", "Bearer "+rule.Http.Auth.Value)
		case "basic":
			req.SetBasicAuth(rule.Http.Auth.Value, "") // Adjust if username/password are needed
		default:
			return fmt.Errorf("unsupported auth type: %s", rule.Http.Auth.Type)
		}
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the response
	fmt.Printf("Response from %s: %s\n", rule.Http.Url, string(body))
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
