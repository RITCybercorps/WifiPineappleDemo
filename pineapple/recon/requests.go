package recon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.ritsec.cloud/BradHacker/ssid-jungle/pineapple"
)

// StartScan starts a recon scan using the WiFi Pineapple Mk7 REST API
func StartScan(token string, live bool, scanTime int, band ReconScanBand) (bool, error) {
	payload, err := json.Marshal(ReconStartBody{
		Live:     live,
		ScanTime: scanTime,
		Band:     band,
	})
	if err != nil {
		return false, fmt.Errorf("failed to marshal login body: %v", err)
	}

	req, err := http.NewRequest("POST", "http://172.16.42.1:1471/api/recon/start", bytes.NewReader(payload))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to POST to /api/recon/start: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var startRes ReconStartResponse
		err = json.Unmarshal(body, &startRes)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/recon/start response: %v", err)
		}
		return startRes.ScanRunning, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/recon/start error object: %v", err)
		}
		return false, fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}

// StopScan stops a recon scan using the WiFi Pineapple Mk7 REST API
func StopScan(token string) (bool, error) {
	req, err := http.NewRequest("POST", "http://172.16.42.1:1471/api/recon/stop", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to POST to /api/recon/stop: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var stopRes ReconStopResponse
		err = json.Unmarshal(body, &stopRes)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/recon/stop response: %v", err)
		}
		return stopRes.Success, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/recon/stop error object: %v", err)
		}
		return false, fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}

// StopScan stops a recon scan using the WiFi Pineapple Mk7 REST API
func Status(token string) (*ReconStatusResponse, error) {
	req, err := http.NewRequest("GET", "http://172.16.42.1:1471/api/recon/status", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to GET to /api/recon/status: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var statusRes ReconStatusResponse
		err = json.Unmarshal(body, &statusRes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal /api/recon/status response: %v", err)
		}
		return &statusRes, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal /api/recon/status error object: %v", err)
		}
		return nil, fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}

// ListScans lists previous recon scans using the WiFi Pineapple Mk7 REST API
func ListScans(token string) ([]ReconListScansResponse, error) {
	req, err := http.NewRequest("GET", "http://172.16.42.1:1471/api/recon/scans", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to GET to /api/recon/scans: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var scansRes []ReconListScansResponse
		err = json.Unmarshal(body, &scansRes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal /api/recon/scans response: %v", err)
		}
		return scansRes, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal /api/recon/scans error object: %v", err)
		}
		return nil, fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}

// DeleteScan deletes a previous recon scans using the WiFi Pineapple Mk7 REST API
func DeleteScan(token string, scanID int) (bool, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://172.16.42.1:1471/api/recon/scans/%d", scanID), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to DELETE to /api/recon/scans: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		var deleteRes ReconDeleteScanResponse
		err = json.Unmarshal(body, &deleteRes)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/recon/scans response: %v", err)
		}
		return deleteRes.Success, nil
	} else {
		var apiError pineapple.APIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal /api/recon/scans error object: %v", err)
		}
		return false, fmt.Errorf("got error message: \"%s\"", apiError.Error)
	}
}
