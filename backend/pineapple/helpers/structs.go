package helpers

type LookupOUIResponse struct {
	Available bool   `json:"available"`
	Vendor    string `json:"vendor"`
}
