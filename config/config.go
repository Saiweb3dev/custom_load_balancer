package config

type Config struct {
	ListenAddr string
	// Add more configuration fields as needed
}

func Load() (*Config, error) {
	// In a real application, you'd load this from a file or environment variables
	return &Config{
		ListenAddr: ":8080",
	}, nil
}