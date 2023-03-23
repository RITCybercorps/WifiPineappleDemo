package pineap

type ListSSIDsResponse struct {
	SSIDs string `json:"ssids"`
}

type ClearSSIDsResponse struct {
	Success bool `json:"success"`
}

type RemoveSSIDBody struct {
	SSID string `json:"ssid"`
}
type RemoveSSIDResponse struct {
	Success bool `json:"success"`
}
