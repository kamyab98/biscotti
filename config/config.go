package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	PixelRedirectKey  string            `mapstructure:"PIXEL_REDIRECT_KEY"`
	CookieKey         string            `mapstructure:"COOKIE_KEY"`
	SchemaRegistryURL string            `mapstructure:"SCHEMA_REGISTRY_URL"`
	NetworkToUrls     map[string]string `mapstructure:"NETWORK_TO_URL"`
	CookieDomain      string            `mapstructure:"COOKIE_DOMAIN"`
	KafkaTopic        string            `mapstructure:"KAFKA_TOPIC"`
}

var appConfig *Config

func GetAppConfig() *Config {
	if appConfig == nil {
		dir, err := os.Getwd()
		config, err := loadConfig(dir)
		if err != nil {
			panic(err)
		}
		appConfig = &config
	}
	return appConfig
}

func loadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	return
}
