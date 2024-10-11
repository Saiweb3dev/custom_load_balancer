package config

// Config holds the configuration for the load balancer
type Config struct {
	ListenAddr string
	// Add more configuration fields as needed
}

// Load retrieves the configuration from a source (e.g., file, environment)
func Load() (*Config, error) {
	// In a real application, you'd load this from a file or environment variables
	return &Config{
		ListenAddr: ":8080",
	}, nil
}