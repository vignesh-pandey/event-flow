package helpers

import (
	"producer-service/constants"
	"producer-service/logs"

	"github.com/spf13/viper"
)

// Read Configuration files
func Configuration(configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(constants.CONFIG_FILE_NAME)

	err := viper.ReadInConfig()
	if err != nil {
		logs.Log.Errorln("Failed to Start Web Server: Viper Read Config Failure -", err)
		panic(err)
	}

	viper.WatchConfig()
}
