package main

import (
	"testing"
)

func TestLoadConfigs(t *testing.T) {
	t.Log("Loading configs")
	res, err := LoadConfigs([]string{
		"./etc/birdwatcher/birdwatcher.conf",
		"/etc/birdwatcher/birdwatcher.conf",
		"./etc/birdwatcher/birdwatcher.local.conf",
	})

	t.Log(res)
	t.Log(err)
}
