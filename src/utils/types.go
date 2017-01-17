package utils

type ArpCache struct {
	Addr   string `json:"addr"`
	HwType string `json:"hw_type"`
	HwAddr string `json:"hw_addr"`
	Mask   string `json:"mask"`
}

type CurrentRam struct {
	Total     float64 `json:"total"`
	Used      float64 `json:"used"`
	Available float64 `json:"available"`
}
