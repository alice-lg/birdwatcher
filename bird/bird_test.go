package bird

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

// Run the test binary as a birdc substitute so these tests do not need BIRD.
func TestRunAndParseCommandHelper(t *testing.T) {
	for i, arg := range os.Args {
		if arg != "--" {
			continue
		}
		args := os.Args[i+1:]
		if len(args) != 3 || args[0] != "-r" || args[1] != "show" {
			t.Fatal("Unexpected birdc arguments:", args)
		}
		if args[2] == "queue-command-error" {
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, "bird output")
		os.Exit(0)
	}
}

func setupRunAndParseTest() func() {
	oldClientConf, oldCache := ClientConf, cache
	RateLimitConf.Lock()
	oldRateLimitConf := RateLimitConf.Conf
	RateLimitConf.Conf = RateLimitConfig{Enabled: false}
	RateLimitConf.Unlock()

	ClientConf = BirdConfig{
		BirdCmd:  os.Args[0] + " -test.run=^TestRunAndParseCommandHelper$ --",
		CacheTtl: 5,
	}
	cache = NewMemoryCache(10)

	return func() {
		ClientConf, cache = oldClientConf, oldCache
		RateLimitConf.Lock()
		RateLimitConf.Conf = oldRateLimitConf
		RateLimitConf.Unlock()
	}
}

type runAndParseTestCache struct {
	Cache
	setError error
	setPanic interface{}
}

func (c runAndParseTestCache) Set(key string, val Parsed, ttl int) error {
	if c.setPanic != nil {
		panic(c.setPanic)
	}
	if c.setError != nil {
		return c.setError
	}
	return c.Cache.Set(key, val, ttl)
}

func TestRunAndParseCleansQueueOnReturn(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		rateLimited bool
		cacheError  bool
		special     Parsed
	}{
		{name: "success", cmd: "queue-success"},
		{name: "rate limited", cmd: "queue-rate-limited", rateLimited: true, special: NilParse},
		{name: "command error", cmd: "queue-command-error", special: BirdError},
		{name: "cache error", cmd: "queue-cache-error", cacheError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup := setupRunAndParseTest()
			defer cleanup()
			defer RunQueue.Delete(test.cmd)
			if test.rateLimited {
				RateLimitConf.Lock()
				RateLimitConf.Conf = RateLimitConfig{Enabled: true, Reqs: 0}
				RateLimitConf.Unlock()
			}
			if test.cacheError {
				cache = runAndParseTestCache{Cache: cache, setError: errors.New("expected cache error")}
			}

			parsed, fromCache := RunAndParse(false, test.cmd, test.cmd, func(io.Reader) Parsed {
				return Parsed{"value": 1}
			}, func(parsed *Parsed) {
				(*parsed)["updated"] = true
			})
			if fromCache {
				t.Error("Expected a fresh result")
			}
			if test.rateLimited || test.cmd == "queue-command-error" {
				if !reflect.DeepEqual(parsed, test.special) {
					t.Errorf("Unexpected error result: %v", parsed)
				}
			} else if parsed["value"] != 1 || parsed["updated"] != true {
				t.Errorf("Unexpected parsed result: %v", parsed)
			}
			if _, ok := RunQueue.Load(test.cmd); ok {
				t.Error("Completed command remained in RunQueue")
			}
		})
	}
}

func TestRunAndParseCachedRequestsLeaveQueueOwner(t *testing.T) {
	for _, useCache := range []bool{true, false} {
		t.Run(fmt.Sprintf("useCache-%t", useCache), func(t *testing.T) {
			cleanup := setupRunAndParseTest()
			defer cleanup()
			cmd := "queue-owned"
			defer RunQueue.Delete(cmd)
			if err := cache.Set(cmd, Parsed{"value": 1}, 5); err != nil {
				t.Fatal(err)
			}

			// Keep a completed owner's entry to exercise both the cache-hit
			// and duplicate-request paths without launching another command.
			group := &sync.WaitGroup{}
			RunQueue.Store(cmd, group)
			parsed, fromCache := RunAndParse(useCache, cmd, cmd, func(io.Reader) Parsed {
				t.Error("Cached request unexpectedly invoked its parser")
				return Parsed{"value": 2}
			}, nil)
			if !fromCache || parsed["value"] != 1 {
				t.Errorf("Expected the cached result, got %v (from cache: %v)", parsed, fromCache)
			}
			if queued, ok := RunQueue.Load(cmd); !ok || queued != group {
				t.Error("Cached request removed another request's queue entry")
			}
		})
	}
}

func TestRunAndParseCleansQueueAfterPanic(t *testing.T) {
	for _, stage := range []string{"parser", "cache update", "cache write"} {
		t.Run(stage, func(t *testing.T) {
			cleanup := setupRunAndParseTest()
			defer cleanup()
			cmd := "queue-panic"
			defer RunQueue.Delete(cmd)
			baseCache := cache
			panicValue := "expected " + stage + " panic"
			if stage == "cache write" {
				cache = runAndParseTestCache{Cache: cache, setPanic: panicValue}
			}

			var group *sync.WaitGroup
			var recovered interface{}
			func() {
				defer func() { recovered = recover() }()
				RunAndParse(false, cmd, cmd, func(io.Reader) Parsed {
					queued, _ := RunQueue.Load(cmd)
					group = queued.(*sync.WaitGroup)
					if stage == "parser" {
						panic(panicValue)
					}
					return Parsed{"value": 1}
				}, func(*Parsed) {
					if stage == "cache update" {
						panic(panicValue)
					}
				})
			}()

			if recovered != panicValue {
				t.Fatalf("Expected panic %q to propagate, got %v", panicValue, recovered)
			}
			if _, ok := RunQueue.Load(cmd); ok {
				t.Error("Panicking command remained in RunQueue")
				RunQueue.Delete(cmd)
			}

			done := make(chan struct{})
			go func() {
				group.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Error("Panicking command did not release its wait group")
				// Release the waiter even when this regression fails.
				group.Done()
				<-done
			}

			cache = baseCache
			parsed, fromCache := RunAndParse(false, cmd, cmd, func(io.Reader) Parsed {
				return Parsed{"value": 2}
			}, nil)
			if fromCache || parsed["value"] != 2 {
				t.Errorf("Expected the same command to run again, got %v (from cache: %v)", parsed, fromCache)
			}
			if _, ok := RunQueue.Load(cmd); ok {
				t.Error("Retried command remained in RunQueue")
			}
		})
	}
}
