package pineap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend/pineapple"
)

// ListSSIDs lists PineAP SSID Pool using the WiFi Pineapple Mk7 REST API
func ListSSIDs(token string) (string, error) {
	req, err := http.NewRequest("GET", "http://172.16.42.1:1471/api/pineap/ssids", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to GET to /api/pineap/ssids: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var listRes ListSSIDsResponse
		err = json.Unmarshal(body, &listRes)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal /api/pineap/ssids response: %v", err)
		}
		return listRes.SSIDs, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal /api/pineap/ssids error object: %v", err)
		}
		return "", fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}

// ClearSSIDs clears PineAP SSID Pool using the WiFi Pineapple Mk7 REST API
func ClearSSIDs(token string) (bool, error) {
	req, err := http.NewRequest("DELETE", "http://172.16.42.1:1471/api/pineap/ssids", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to DELETE to /api/pineap/ssids: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var clearRes ClearSSIDsResponse
		err = json.Unmarshal(body, &clearRes)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/pineap/ssids response: %v", err)
		}
		return clearRes.Success, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/pineap/ssids error object: %v", err)
		}
		return false, fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}

// RemoveSSID removes a single SSID from the PineAP SSID Pool using the WiFi Pineapple Mk7 REST API
func RemoveSSID(token, ssid string) (bool, error) {
	payload, err := json.Marshal(RemoveSSIDBody{
		SSID: ssid,
	})
	if err != nil {
		return false, fmt.Errorf("failed to marshal login body: %v", err)
	}

	req, err := http.NewRequest("DELETE", "http://172.16.42.1:1471/api/pineap/ssids/ssid", bytes.NewReader(payload))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to DELETE to /api/pineap/ssids/ssid: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var removeRes RemoveSSIDResponse
		err = json.Unmarshal(body, &removeRes)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/pineap/ssids/ssid response: %v", err)
		}
		return removeRes.Success, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/pineap/ssids/ssid error object: %v", err)
		}
		return false, fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}
