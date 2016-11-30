package main

// Birdwatcher Configuration

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

type Config struct {
	Server ServerConfig
	Status StatusConfig
	Bird   BirdConfig
	Bird6  BirdConfig
}

type ServerConfig struct {
	AllowFrom []string `toml:"allow_from"`
}

type StatusConfig struct {
	ReconfigTimestampSource string `toml:"reconfig_timestamp_source"`
	ReconfigTimestampMatch  string `toml:"reconfig_timestamp_match"`

	FilteredFields []string `toml:"filtered_fields"`
}

type BirdConfig struct {
	Listen         string
	ConfigFilename string `toml:"config"`
	BirdCmd        string `toml:"birdc"`
}

// Try to load configfiles as specified in the files
// list. For example:
//
//    ./etc/birdwatcher/birdwatcher.conf
//    /etc/birdwatcher/birdwatcher.conf
//    ./etc/birdwatcher/birdwatcher.local.conf
//
//
func LoadConfigs(configFiles []string) (*Config, error) {
	config := &Config{}
	hasConfig := false
	var confError error

	for _, filename := range configFiles {
		tmp := &Config{}
		_, err := toml.DecodeFile(filename, tmp)
		if err != nil {
			continue
		} else {
			hasConfig = true
			// Merge configs
			if err := mergo.Merge(config, tmp); err != nil {
				return nil, err
			}
		}
	}

	if !hasConfig {
		confError = fmt.Errorf("Could not load any config file")
	}

	return config, confError
}
