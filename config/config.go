package config

import (
	"github.com/spf13/viper"
)

type DbConfig struct {
	Url      string `mapstructure:"url"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
}

type RiotConfig struct {
	Key            string  `mapstructure:"key"`
	RateLimit      int     `mapstructure:"rate_limit"`
	RatePeriod     int     `mapstructure:"rate_period"`
	RateEfficiency float32 `mapstructure:"rate_efficiency"`
	MatchesAfter   int64   `mapstructure:"matches_after"`
}

type Config struct {
	Db   DbConfig   `mapstructure:"db"`
	Riot RiotConfig `mapstructure:"riot"`
}

func LoadConfig() (Config, error) {
	vp := viper.New()
	vp.SetConfigFile("./config/config.json")

	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
