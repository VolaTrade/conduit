package models

type CacheEntry struct {
	TxUrl string
	ObUrl string
	Pair  string
}

type CortexEntry struct {
	Url string
}

type CortexRequest struct {
	Data string `json:"data"`
}
