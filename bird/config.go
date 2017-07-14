package bird

// Birdwatcher Birdc Configuration

type StatusConfig struct {
	ReconfigTimestampSource string `toml:"reconfig_timestamp_source"`
	ReconfigTimestampMatch  string `toml:"reconfig_timestamp_match"`

	FilterFields []string `toml:"filter_fields"`
}

type BirdConfig struct {
	Listen         string
	ConfigFilename string `toml:"config"`
	BirdCmd        string `toml:"birdc"`
}

type ParserConfig struct {
	FilterFields []string `toml:"filter_fields"`

	PerPeerTables        bool   `toml:"per_peer_tables"`
	PeerProtocolMatch    string `toml:"peer_protocol_match"`
	PipeProtocolTemplate string `toml:"pipe_protocol_template"`
}

type RateLimitConfig struct {
	Reqs    int
	Max     int `toml:"requests_per_minute"`
	Enabled bool
}
