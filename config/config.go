package config

import "github.com/spf13/viper"

type Conf struct {
	ServerPort             int    `mapstructure:"SERVER_PORT"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              int    `mapstructure:"REDIS_PORT"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	RedisDB                int    `mapstructure:"REDIS_DB"`
	IPMaxRequests          int    `mapstructure:"IP_MAX_REQUESTS"`
	TokenMaxRequests       int    `mapstructure:"TOKEN_MAX_REQUESTS"`
	TimeWindowMilliseconds int    `mapstructure:"TIME_WINDOW_MILISECONDS"`
}

func Load(path string) (*Conf, error) {
	var c *Conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	return c, nil
}
