package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const Version = "0.1.0"

const dbFileName = "social.db"

var ServerURI string

var (
	PrivKey  string
	PubKey   string
	PathToDB string
)

func LoadFrom(socialRoot string) {
	insideDocker := os.Getenv("INSIDE_DOCKER") != ""

	setupDefaultConfig(insideDocker)

	pathToConfig := filepath.Join(socialRoot, "config.toml")
	if _, err := os.Stat(pathToConfig); err == nil {
		viper.SetConfigFile(pathToConfig)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("failed to load configuration file %s:\n\t%s\n", pathToConfig, err)
		}
		log.Printf("Successfully loaded configuration (config file = %s)", pathToConfig)
	} else {
		log.Printf("Warning: no getascension.toml config file found in directory %s\n", socialRoot)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("_", ""))
	viper.SetEnvPrefix("getascension")
}
