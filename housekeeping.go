package main

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/alice-lg/birdwatcher/bird"
)

type HousekeepingConfig struct {
	Interval           int  `toml:"interval"`
	ForceReleaseMemory bool `toml:"force_release_memory"`
}

// This is used to run regular housekeeping tasks, currently expiring old
// Cache entries to release memory
func Housekeeping(config HousekeepingConfig, expireCaches bool) {
	for {
		if config.Interval > 0 {
			time.Sleep(time.Duration(config.Interval) * time.Minute)
		} else {
			time.Sleep(5 * time.Minute)
		}

		log.Println("Housekeeping started")

		if (bird.ClientConf.CacheTtl > 0) && expireCaches {
			// Expire the caches
			log.Println("Expiring MemoryCache")

			count := bird.ExpireCache()
			log.Println("Expired", count, "entries (MemoryCache)")
		}

		if config.ForceReleaseMemory {
			// Trigger a GC and SCVG run
			log.Println("Freeing memory")
			debug.FreeOSMemory()
		}
	}
}
