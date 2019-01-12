package config

import (
	"github.com/spf13/viper"
)

func setupDefaultConfig(insideDocker bool) {
	viper.SetTypeByDefaultValue(true)

	viper.SetDefault("InsideDocker", insideDocker)

	viper.SetDefault("ServerName", "getascension.com")
	viper.SetDefault("ServerTitle", "Ascension")
	viper.SetDefault("ServerBaseURL", "https://getascension.com/")
	viper.SetDefault("ServerFQDN", "getascension.com")

	viper.SetDefault("ImageRoot", "/images")

	viper.SetDefault("ServerValidDNS", map[string]interface{}{
		"getascension.com": true,
	})
}
