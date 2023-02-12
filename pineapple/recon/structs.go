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
