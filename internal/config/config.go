package config

import "github.com/spf13/viper"

type Config struct {
	Store string `json:"store,omitempty"`
	Redis Redis  `json:"redis,omitempty"`
}

type Redis struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
	Database int    `json:"database,omitempty"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}
	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return cfg, nil
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	err := viper.Unmarshal(cfg)
	return cfg, err
}
