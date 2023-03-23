package recon

const (
	ReconCONTINUOUS int = 0
)

type ReconScanBand string

const (
	Recon2_4GHZ      ReconScanBand = "0"
	Recon5GHZ        ReconScanBand = "1"
	Recon2_4GHZ_5GHZ ReconScanBand = "2"
)

type ReconBoolInt int

const (
	ReconFalse ReconBoolInt = 0
	ReconTrue  ReconBoolInt = 1
)

type ReconAPExt struct {
	ReconAP
	Vendor string `json:"vendor,omitempty"`
}

type ReconAP struct {
	ScanId        int           `json:"scan_id"`
	Ssid          string        `json:"ssid"`
	Bssid         string        `json:"bssid"`
	Encryption    int           `json:"encryption"`
	Hidden        ReconBoolInt  `json:"hidden"`
	Wps           ReconBoolInt  `json:"wps"`
	Channel       int           `json:"channel"`
	Signal        int           `json:"signal"`
	Data          int           `json:"data"`
	FirstSeen     int64         `json:"first_seen"`
	LastSeen      int64         `json:"last_seen"`
	LastSeenDelta int           `json:"last_seen_delta"`
	Probes        int           `json:"probes"`
	Mfp           ReconBoolInt  `json:"mfp"`
	Clients       []ReconClient `json:"clients"`
	NumClients    int           `json:"num_clients"`
}

type ReconClient struct {
	ScanId          int    `json:"ScanID"`
	ClientMAC       string `json:"client_mac"`
	ApMAC           string `json:"ap_mac"`
	ApChannel       int    `json:"ap_channel"`
	Data            int    `json:"data"`
	BroadcastProbes int    `json:"broadcast_probes"`
	FirstSeen       int64  `json:"first_seen"`
	LastSeen        int64  `json:"last_seen"`
	LastSeenDelta   int    `json:"last_seen_delta"`
	Ssid            string `json:"ssid"`
}

type ReconScanApBody struct {
	Search   string `json:"search"`
	Page     int    `json:"page"`
	PageSize uint   `json:"pageSize"`
	SortCol  string `json:"sortCol"`
	SortDir  string `json:"sortDir"`
}

type ReconScanApResponse struct {
	TotalLength int       `json:"totalLength"`
	Aps         []ReconAP `json:"aps"`
	Page        int       `json:"page"`
}

type ReconStartBody struct {
	Live     bool          `json:"live"`
	ScanTime int           `json:"scan_time"` // In Minutes (set to 0 for continuous)
	Band     ReconScanBand `json:"band"`      // 0 - 2.4 GHz | 1 - 5 GHz | 2 - Both
}

type ReconStartResponse struct {
	ScanRunning bool `json:"scanRunning"`
	ScanID      int  `json:"scanID"`
}

type ReconStopResponse struct {
	Success bool `json:"success"`
}

type ReconStatusResponse struct {
	CaptureRunning bool `json:"captureRunning"`
	ScanRunning    bool `json:"scanRunning"`
	Continuous     bool `json:"continuous"`
	ScanPercent    int  `json:"scanPercent"`
	ScanID         int  `json:"scanID"`
}

type ReconListScansResponse struct {
	ScanID int    `json:"scan_id"`
	Date   string `json:"date"`
}

type ReconDeleteScanResponse struct {
	Success bool `json:"success"`
}
