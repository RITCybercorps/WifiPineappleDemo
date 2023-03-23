package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/backend/pineapple"
)

// ListSSIDs lists PineAP SSID Pool using the WiFi Pineapple Mk7 REST API
func LookupOUI(token string, oui string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://172.16.42.1:1471/api/helpers/lookupOUI/%s", oui), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to GET to /api/helpers/lookupOUI/%s: %v", oui, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var lookupRes LookupOUIResponse
		err = json.Unmarshal(body, &lookupRes)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal /api/helpers/lookupOUI/%s response: %v", oui, err)
		}
		return lookupRes.Vendor, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal /api/helpers/lookupOUI/%s error object: %v", oui, err)
		}
		return "", fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}
