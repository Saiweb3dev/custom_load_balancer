package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the configuration for the load balancer
type Config struct {
	// Server settings
	ListenAddr string `json:"listen_addr"`

	MongoURI string `json:"mongo_uri"`
    MongoDB  string `json:"mongo_db"`
	
	// TLS settings
	TLSCertFile string `json:"tls_cert_file"`
	TLSKeyFile  string `json:"tls_key_file"`
	
	// Connection pool settings
	PoolMaxConns        int           `json:"pool_max_conns"`
	PoolIdleTimeout     time.Duration `json:"pool_idle_timeout"`
	PoolMaxLifetime     time.Duration `json:"pool_max_lifetime"`
	PoolCleanupInterval time.Duration `json:"pool_cleanup_interval"`
	
	// Health checker settings
	HealthCheckInterval time.Duration `json:"health_check_interval"`
	HealthCheckTimeout  time.Duration `json:"health_check_timeout"`
	HealthCheckEndpoint string        `json:"health_check_endpoint"`
	
	// Registry settings
	RegistryFile string `json:"registry_file"`
	
	// Balancer settings
	BalancerAlgorithm string `json:"balancer_algorithm"`
	
	// Logging settings
	LogLevel  string `json:"log_level"`
	LogFormat string `json:"log_format"`
}

// Load retrieves the configuration from a JSON file
func Load(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	// Set default values for fields that are not specified
	config.setDefaults()

	return config, nil
}

// setDefaults sets default values for configuration fields that are not specified
func (c *Config) setDefaults() {
	if c.ListenAddr == "" {
		c.ListenAddr = ":8080"
	}

	if c.MongoURI == "" {
		c.MongoURI = "mongodb://localhost:27017"
}
if c.MongoDB == "" {
		c.MongoDB = "userdb"
}
	if c.PoolMaxConns == 0 {
		c.PoolMaxConns = 100
	}
	if c.PoolIdleTimeout == 0 {
		c.PoolIdleTimeout = 5 * time.Minute
	}
	if c.PoolMaxLifetime == 0 {
		c.PoolMaxLifetime = 30 * time.Minute
	}
	if c.PoolCleanupInterval == 0 {
		c.PoolCleanupInterval = 1 * time.Minute
	}
	if c.HealthCheckInterval == 0 {
		c.HealthCheckInterval = 10 * time.Second
	}
	if c.HealthCheckTimeout == 0 {
		c.HealthCheckTimeout = 5 * time.Second
	}
	if c.HealthCheckEndpoint == "" {
		c.HealthCheckEndpoint = "/health"
	}
	if c.RegistryFile == "" {
		c.RegistryFile = "registry.json"
	}
	if c.BalancerAlgorithm == "" {
		c.BalancerAlgorithm = "round_robin"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.LogFormat == "" {
		c.LogFormat = "text"
	}
}