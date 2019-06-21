package config

// Config is the config struct for the application.
type Config struct {
	Clients struct {
		Jokes struct {
			Name    string `mapstructure:"name"`
			BaseURL string `mapstructure:"base_url"`
			Timeout int    `mapstructure:"timeout"`
		} `mapstructure:"joke_service"`
		Names struct {
			Name    string `mapstructure:"name"`
			BaseURL string `mapstructure:"base_url"`
			Timeout int    `mapstructure:"timeout"`
		} `mapstructure:"name_service"`
	} `mapstructure:"clients"`
}
