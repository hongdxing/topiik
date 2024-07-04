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

type ServerConfig struct {
	Listen string
	Role   string // controller, broker
	//Join       string // comma seprated host list
	SaveMillis uint // Persistent Job interval

	/*** Internal Use Only***/
	RaftHeartbeatMin uint16 // Raft random heartbeat Min
	RaftHeartbeatMax uint16 // Raft random heartbeat Max
	//JoinList         []string // Internal use
	//NodeRole uint8 // Internal use

	Host    string
	Port    string
	CA_PORT string
}

func ParseServerConfig(configPath string) (*ServerConfig, error) {
	serverConfig := ServerConfig{
		Host:    "localhost",
		Port:    "8301",
		Listen:  "localhost:8301",
		CA_PORT: "18301",
	}
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

	err := viper.Unmarshal(&serverConfig)
	if err != nil {
		fmt.Println(err)
	}
	/*
		serverConfig.Join = strings.Trim(serverConfig.Join, consts.SPACE)
		if serverConfig.Join != "" {
			serverConfig.JoinList = strings.Split(serverConfig.Join, ",")
			for i, s := range serverConfig.JoinList {
				serverConfig.JoinList[i] = strings.Trim(s, consts.SPACE)
			}
		}*/

	fmt.Printf("Using config file: %s\n", theConfigPath)

	// set CA_PORT
	reg := regexp.MustCompile(`(.*)((?::))((?:[0-9]+))$`)
	pieces := reg.FindStringSubmatch(serverConfig.Listen)
	if len(pieces) != 4 {
		return nil, errors.New("Invalid Listen format: " + serverConfig.Listen)
	}
	iPort, err := strconv.Atoi(pieces[3])
	if err != nil {
		return nil, errors.New("Invalid Listen format: " + serverConfig.Listen)
	}
	serverConfig.CA_PORT = strconv.Itoa(10000 + iPort)
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
