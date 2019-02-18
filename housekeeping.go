package main

import (
	"time"
	"log"
	"runtime/debug"

	"github.com/alice-lg/birdwatcher/bird"
)


type HousekeepingConfig struct {
	ForceReleaseMemory bool `toml:"force_release_memory"`
}

// This is used to run regular housekeeping tasks, currently expiring old
// Cache entries to release memory
func Housekeeping(config HousekeepingConfig) {
	for {
		if bird.ClientConf.CacheTtl > 0 {
			time.Sleep(time.Duration(bird.ClientConf.CacheTtl) * time.Minute)
		} else {
			time.Sleep(5 * time.Minute)
		}

		log.Println("Housekeeping started")

		if bird.ClientConf.CacheTtl > 0 {
			// Expire the caches
			log.Println("Expiring caches")

			count := bird.ParsedCache.Expire()
			log.Println("Expired", count, "entries (ParsedCache)")

			count = bird.MetaCache.Expire()
			log.Println("Expired", count, "entries (MetaCache)")
		}

		if config.ForceReleaseMemory {
			// Trigger a GC and SCVG run
			log.Println("Freeing memory")
			debug.FreeOSMemory()
		}
	}
}
