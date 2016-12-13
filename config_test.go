package main

import (
	"testing"
)

func TestLoadConfigs(t *testing.T) {
	t.Log("Loading configs")
	res, err := LoadConfigs([]string{
		"./etc/ecix/birdwatcher.conf",
		"/etc/ecix/birdwatcher.conf",
		"./etc/ecix/birdwatcher.local.conf",
	})

	t.Log(res)
	t.Log(err)
}
