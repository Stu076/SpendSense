package config

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
)

const (
	configDir     = "./config/"
	configFile    = ".env"
	configDefault = ".env.default"
)

func Init() bool {
	isDone, err := checkForConfig()

	if err != nil {
		fmt.Println("Something went wrong initializing the config: ", err)
		return false
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigType("env")
	viper.SetConfigName(configFile)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		fmt.Println("Error reading config file: ", err)
		return false
	}

	return isDone
}

func checkForConfig() (bool, error) {
	fileInfo, err := os.Stat(configDir + "/" + configFile)
	if err == nil {
		return !fileInfo.IsDir(), nil
	}

	err = createConfig()
	return err != nil, err
}

func createConfig() error {
	defaultConf, err := os.Open(configDir + "/" + configDefault)
	if err != nil {
		return err
	}

	defer defaultConf.Close()

	conf, err := os.Create(configDir + "/" + configFile)

	if err != nil {
		return err
	}

	defer conf.Close()

	_, err = io.Copy(conf, defaultConf)

	if err != nil {
		return err
	}

	err = conf.Sync()
	return err
}
