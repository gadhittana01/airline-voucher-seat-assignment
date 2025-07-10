package utils

import (
	"time"

	"github.com/spf13/viper"
)

type BaseConfig struct {
	Environment         string        `mapstructure:"ENVIRONMENT"`
	ServerPort          int           `mapstructure:"SERVER_PORT"`
	DBPath              string        `mapstructure:"DB_PATH"`
	LogLevel            string        `mapstructure:"LOG_LEVEL"`
	CORSAllowedOrigins  []string      `mapstructure:"CORS_ALLOWED_ORIGINS"`
	JWTKey              string        `mapstructure:"JWT_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadBaseConfig(path string, configName string, config *BaseConfig) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func CheckAndSetConfig(path string, configName string) *BaseConfig {
	config := &BaseConfig{}
	LoadBaseConfig(path, configName, config)
	return config
}
