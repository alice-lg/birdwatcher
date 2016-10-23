package main

type TimeInfo struct {
	Date         string `json:"date"`
	TimezoneType string `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type CacheStatus struct {
	OrigTTL  int      `json:"orig_ttl"`
	CachedAt TimeInfo `json:"cached_at"`
}

type APIInfo struct {
	Version         string
	ResultFromCache bool        `json:"result_from_cache"`
	CacheStatus     CacheStatus `json:"cache_status"`
}

func GetApiInfo() *APIInfo {
	ai := &APIInfo{}

	/* Dummy data until we implement caching */
	ai.Version = "1.0"
	ai.ResultFromCache = false

	return ai
}
