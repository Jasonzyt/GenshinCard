package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/viper"
)

const ConfigFilePath = "./config.yaml"

type Config struct {
	CookieSecret string `default:""`
	HttpServer   struct {
		Host string `default:"0.0.0.0"`
		Port int    `default:"8080"`
		Ssl  struct {
			Enabled  bool   `default:"false"`
			CertFile string `default:"cert.pem"`
			KeyFile  string `default:"key.pem"`
		}
		AccessKey string `default:""`
	}
}

func setDefaultValues(v interface{}, key string) {
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			setDefaultValues(reflect.New(field.Type).Elem().Interface(), key+field.Name+".")
		} else {
			defaultString := field.Tag.Get("default")
			var defaultValue interface{}
			switch field.Type.Kind() {
			case reflect.String:
				defaultValue = defaultString
			case reflect.Int:
				defaultValue, _ = strconv.Atoi(defaultString)
			case reflect.Bool:
				defaultValue, _ = strconv.ParseBool(defaultString)
			case reflect.Float64:
				defaultValue, _ = strconv.ParseFloat(defaultString, 64)
			}
			viper.SetDefault(key+field.Name, defaultValue)
		}
	}
}

func initConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	setDefaultValues(Config{}, "")
	file, err := os.Open(ConfigFilePath)
	file.Close()
	if err != nil {
		fmt.Println("[WARN] config.yaml not found, creating...")
		viper.WriteConfigAs(ConfigFilePath)
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
	viper.WriteConfigAs(ConfigFilePath)
	return &config
}

var GlobalConfig = &Config{}

func Init() {
	GlobalConfig = initConfig()
}
