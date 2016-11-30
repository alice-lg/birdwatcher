package endpoints

// Endpoints / Server configuration
type ServerConfig struct {
	AllowFrom []string `toml:"allow_from"`
}
