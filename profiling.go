package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// Write a heap profile to the given file.
func createHeapProfile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

// Write a memory allocation profile to the given file.
func createAllocProfile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create alloc profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
		log.Fatal("could not write alloc profile: ", err)
	}
}

// Start a goroutine to periodically write memory profiles.
func startMemoryProfile(prefix string) {
	t := 0
	log.Println("Starting memory profiling:", prefix)
	for {
		filename := fmt.Sprintf("%s-heap-%03d", prefix, t)
		runtime.GC() // get up-to-date statistics (according to docs)
		createHeapProfile(filename)
		log.Println("Wrote memory heap profile:", filename)
		filename = fmt.Sprintf("%s-allocs-%03d", prefix, t)
		log.Println("Wrote memory allocs profile:", filename)
		createAllocProfile(filename)
		time.Sleep(30 * time.Second)
		t++
	}
}
