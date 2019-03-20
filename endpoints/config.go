package endpoints

// Endpoints / Server configuration
type ServerConfig struct {
	AllowFrom      []string `toml:"allow_from"`
	ModulesEnabled []string `toml:"modules_enabled"`
	AllowUncached  bool     `toml:"allow_uncached"`

	EnableTLS bool   `toml:"enable_tls"`
	Crt       string `toml:"crt"`
	Key       string `toml:"key"`
}
