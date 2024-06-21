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
	"strings"
	"topiik/internal/consts"
	"topiik/raft"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host     string
	Port     string
	Listen   string
	Join     string // comma seprated host list
	JoinList []string
	NodeRole uint8
}

func ParseServerConfig(configPath string) *ServerConfig {
	serverConfig := ServerConfig{
		Host:   "localhost",
		Port:   "8302",
		Listen: "localhost:8302",
		Join:   "",
	}
	theConfigPath := "server.env"
	if configPath != "" {
		_, error := os.Stat(configPath)
		if error != nil {
			fmt.Printf("config file %s not exists\n", configPath)
			return &ServerConfig{}
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
	serverConfig.Join = strings.Trim(serverConfig.Join, consts.SPACE)
	if serverConfig.Join != "" {
		serverConfig.JoinList = strings.Split(serverConfig.Join, ",")
		for i, s := range serverConfig.JoinList {
			serverConfig.JoinList[i] = strings.Trim(s, consts.SPACE)
		}
	}
	fmt.Printf("Using config file: %s\n", theConfigPath)
	printConfig(serverConfig)
	// when server start, default to FOLLOWER
	serverConfig.NodeRole = raft.NODE_LEADER
	return &serverConfig
}

func printConfig(serverConfig ServerConfig) {
	fmt.Println(serverConfig)
}
