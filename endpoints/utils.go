package endpoints

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

// go generate does not work in subdirectories. Beautious.
var VERSION string

func GetApiInfo(from_cache bool) *APIInfo {
	ai := &APIInfo{}

	ai.Version = VERSION
	ai.ResultFromCache = from_cache

	return ai
}
