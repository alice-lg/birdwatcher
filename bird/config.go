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
	CacheTtl       int    `toml:"ttl"`
	Dualstack      bool   `toml:"dualstack"`
}

type ParserConfig struct {
	FilterFields []string `toml:"filter_fields"`
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

	MaxKeys int `toml:"max_keys"`
}
