package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func ReadConfig() *Config {
	viper.SetConfigName("app-config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("KTT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	return &cfg
}

type (
	Config struct {
		Env      string
		Port     int
		DB       DBConfig
		CronTabs CronTabs
		Banks    Banks
	}

	DBConfig struct {
		User     string
		Password string
		Host     string
		Port     int
		Name     string
	}

	CronTabs struct {
		CheckOffersCronTab string
	}

	Banks struct {
		FastBankURL  string
		SolidBankURL string
	}
)
