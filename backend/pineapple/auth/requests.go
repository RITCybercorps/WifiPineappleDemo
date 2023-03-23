package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Login generates a user session token for the WiFi Pineapple Mk7 REST API
func Login(username, password string) (string, error) {
	payload, err := json.Marshal(LoginBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal login body: %v", err)
	}

	res, err := http.Post("http://172.16.42.1:1471/api/login", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to POST to /api/login: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var loginRes LoginResponse
	err = json.Unmarshal(body, &loginRes)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal /api/login response: %v", err)
	}

	return loginRes.Token, nil
}
