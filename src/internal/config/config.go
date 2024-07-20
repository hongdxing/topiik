/*
* author: duan hongxing
* date: 21 Jun 2024
* desc: Server configuration types
*
 */

package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/viper"
)

type envConfig struct {
	Listen     string
	SaveMillis uint // Persistent Job interval
}

type ServerConfig struct {
	Listen     string
	SaveMillis uint // Persistent Job interval

	/*** Internal Use Only***/
	RaftHeartbeatMin uint16 // Raft random heartbeat Min
	RaftHeartbeatMax uint16 // Raft random heartbeat Max
	//JoinList         []string // Internal use
	//NodeRole uint8 // Internal use

	Host  string
	Port  string
	PORT2 string // for Cluster use
}

func ParseServerConfig(configPath string) (*ServerConfig, error) {
	config := envConfig{
		Listen: "localhost:8301",
	}
	serverConfig := ServerConfig{}
	theConfigPath := "server.env"
	if configPath != "" {
		_, error := os.Stat(configPath)
		if error != nil {
			fmt.Printf("config file %s not exists\n", configPath)
			return nil, errors.New("config file not exist: " + configPath)
		}
		theConfigPath = configPath
	}
	viper.SetConfigFile(theConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file %s, %s\n", theConfigPath, err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Using config file: %s\n", theConfigPath)

	serverConfig.Listen = config.Listen
	serverConfig.SaveMillis = config.SaveMillis

	// set PORT2 to PORT + 10000
	reg := regexp.MustCompile(`(.*)((?::))((?:[0-9]+))$`)
	pieces := reg.FindStringSubmatch(serverConfig.Listen)
	if len(pieces) != 4 {
		return nil, errors.New("Invalid Listen format: " + serverConfig.Listen)
	}
	iPort, err := strconv.Atoi(pieces[3])
	if err != nil {
		return nil, errors.New("Invalid Listen format: " + serverConfig.Listen)
	}
	serverConfig.PORT2 = strconv.Itoa(10000 + iPort)
	serverConfig.Host = pieces[1]
	serverConfig.Port = pieces[3] // pieces[1] is ":"

	// when server start, default to FOLLOWER
	//serverConfig.NodeRole = raft.ROLE_LEADER

	// Set Raft Heartbeat
	serverConfig.RaftHeartbeatMin = 300
	serverConfig.RaftHeartbeatMax = 600

	// set default SaveMillis
	if serverConfig.SaveMillis == 0 {
		serverConfig.SaveMillis = 1000
	}
	printConfig(serverConfig)
	return &serverConfig, nil
}

func printConfig(serverConfig ServerConfig) {
	fmt.Println(serverConfig)
}
