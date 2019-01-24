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
	FilterFields       []string `toml:"filter_fields"`
	PerPeerTables      bool     `toml:"per_peer_tables"`
	PeerProtocolPrefix string   `toml:"peer_protocol_prefix"`
	PipeProtocolPrefix string   `toml:"pipe_protocol_prefix"`
}

type RateLimitConfig struct {
	Reqs    int
	Max     int `toml:"requests_per_minute"`
	Enabled bool
}

type CacheConfig struct {
	UseRedis      bool   `toml:"use_redis"`
	RedisServer   string `toml:"redis_server"`
	RedisPassword string `toml:"redis_password"`
	RedisDb       int    `toml:"redis_db"`
}
