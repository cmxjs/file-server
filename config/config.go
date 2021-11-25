package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	Port      string
	Cache     string
	RedisAddr string
	RedisDB   string
)

// create config file
func createDefaultConfig(configFile string, defaultCache string) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yml")

	viper.SetDefault("port", "8080")
	viper.SetDefault("cache", defaultCache)
	viper.SetDefault("redis_addr", "localhost:6379")
	viper.SetDefault("redis_db", "0")

	if err := viper.WriteConfig(); err != nil {
		log.Fatalln(err)
	}
}

// read config from config file
func getConfigFromConfigFile(configFile string) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	Port = viper.GetString("port")
	Cache = viper.GetString("cache")
	RedisAddr = viper.GetString("redis_addr")
	RedisDB = viper.GetString("redis_db")

	if Port == "" {
		log.Fatalln("Failed to get field 'port' from config file")
	}
	if Cache == "" {
		log.Fatalln("Failed to get field 'cache' from config file")
	}
}

func init() {
	absDir := func() string {
		absPath, err := filepath.Abs(os.Args[0])
		if err != nil {
			panic(err)
		}
		return filepath.Dir(absPath)
	}()

	defaultCache := filepath.Join(absDir, "cache")
	if _, err := os.Stat(defaultCache); err != nil {
		if err := os.MkdirAll(defaultCache, 0755); err != nil { // create when cache not exists
			log.Fatalf("Failed to create cache dir. cache dir: %v\n", defaultCache)
		}
	}

	configFile := filepath.Join(absDir, "config.yml")
	if _, err := os.Stat(configFile); err != nil {
		createDefaultConfig(configFile, defaultCache)
	}
	getConfigFromConfigFile(configFile)
}
