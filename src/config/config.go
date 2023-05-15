package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	CookieSecret string
}

func updateConfig() {
	viper.Set("CookieSecret", GlobalConfig.CookieSecret)
}

func initConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	file, err := os.Open(filepath.Join(".", "config.yaml"))
	file.Close()
	if err != nil {
		fmt.Println("[WARN] config.yaml not found, creating...")
		updateConfig()
		viper.WriteConfigAs(filepath.Join(".", "config.yaml"))
	}
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("[ERROR] read config failed: %v\n", err)
		return nil
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("[ERROR] unmarshal config failed: %v\n", err)
		return nil
	}
	return &config
}

var GlobalConfig = &Config{}

func Init() {
	GlobalConfig = initConfig()
}
