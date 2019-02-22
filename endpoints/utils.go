package endpoints

import (
	"time"

	"github.com/alice-lg/birdwatcher/bird"
)

type TimeInfo struct {
	Date         time.Time `json:"date"`
	TimezoneType string    `json:"timezone_type"`
	Timezone     string    `json:"timezone"`
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

func GetApiInfo(res *bird.Parsed, from_cache bool) *APIInfo {
	ai := &APIInfo{}

	ai.Version = VERSION
	ai.ResultFromCache = from_cache

	api := *res

	// Derive cache status from TTL
	cachedAt, ok := api["cached_at"].(time.Time)
	if !ok {
		cachedAt = time.Time{}
	}

	// tbh. I have no clue what the difference between
	// timezone type and timezone actually is.
	// I could trace back the timezonetype to the symphony framework
	// Barry was using; the docs say it accepts timezones like
	// "America/New_York", however nothing about UTC could be found.
	//
	// As we convert everything to UTC and let the client
	// render it in local time, it is safe to set this to a fixed
	// value.

	cacheInfo := CacheStatus{
		OrigTTL: -1,
		CachedAt: TimeInfo{
			Date:         cachedAt,
			TimezoneType: "UTC",
			Timezone:     "UTC",
		},
	}

	ai.CacheStatus = cacheInfo

	return ai
}
