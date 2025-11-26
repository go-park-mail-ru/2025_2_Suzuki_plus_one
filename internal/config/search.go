package config

// TODO: use

type SearchServiceConfig struct {
	Host string
	Port string
}


func LoadSearchServiceConfig() SearchServiceConfig {
	return SearchServiceConfig{
		Host: getEnv("SEARCH_SERVICE_HOST", "localhost"),
		Port: getEnv("SEARCH_SERVICE_PORT", "8082"),
	}
}