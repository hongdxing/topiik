/*
* author: duan hongxing
* date: 21 Jun 2024
* desc: Server configuration types
*
 */

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port string
}

func ParseServerConfig(configPath string) ServerConfig {
	serverConfig := ServerConfig{Host: "localhost", Port: "8302"}
	theConfigPath := "server.env"
	if configPath != "" {
		_, error := os.Stat(configPath)
		if error != nil {
			fmt.Printf("config file %s not exists\n", configPath)
			return ServerConfig{}
		}
		theConfigPath = configPath
	}
	viper.SetConfigFile(theConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file %s, %s\n", theConfigPath, err)
	}

	err := viper.Unmarshal(&serverConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Using config file: %s\n", theConfigPath)
	printConfig(serverConfig)

	return serverConfig
}

func printConfig(serverConfig ServerConfig) {
	fmt.Println(serverConfig)
}
