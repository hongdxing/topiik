//author: duan hongxing
//date: 21 Jun 2024
//desc: Server configuration types

package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"topiik/internal/logger"
	"topiik/resp"

	"github.com/spf13/viper"
)

var l = logger.Get()

// node role
const (
	ROLE_WORKER    string = "worker"
	ROLE_PERSISTOR string = "persistor"
)

type envConfig struct {
	Role       string
	Listen     string // current node listen address
	Persistors string // set persistors if current node is controller
	SaveMillis uint   // Persistent Job interval
}

type ServerConfig struct {
	Role       string
	Listen     string
	Persistors []string
	SaveMillis uint // Persistent Job interval

	/*** Internal Use Only***/
	RaftHeartbeatMin uint16 // Raft random heartbeat Min
	RaftHeartbeatMax uint16 // Raft random heartbeat Max
	//JoinList         []string // Internal use
	//NodeRole uint8 // Internal use

	Host  string
	Port  string
	Port2 string // for Cluster use
}

func ParseServerConfig(configPath string) (*ServerConfig, error) {
	config := envConfig{
		//Listen: "localhost:8301",
	}
	serverConfig := ServerConfig{}
	theConfigPath := "server.conf"
	if configPath != "" {
		_, err := os.Stat(configPath)
		if err != nil {
			l.Err(err).Msgf("config file %s not exists\n", configPath)
			return nil, errors.New("config file not exist: " + configPath)
		}
		theConfigPath = configPath
	}

	viper.SetConfigType("env")
	replacer := strings.NewReplacer(".", "")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigFile(theConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		l.Err(err).Msgf("Error reading config file %s, %s", theConfigPath, err)
	}

	// parse config
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)

	l.Info().Msgf("Using config file: %s", theConfigPath)

	serverConfig.Role = config.Role
	serverConfig.Listen = config.Listen
	serverConfig.Persistors = strings.Split(config.Persistors, ",")
	serverConfig.SaveMillis = config.SaveMillis

	// validate role
	if strings.TrimSpace(serverConfig.Role) == "" ||
		(strings.ToLower(serverConfig.Role) != ROLE_WORKER && strings.ToLower(serverConfig.Role) != ROLE_PERSISTOR) {
		return nil, errors.New(resp.RES_SYNTAX_ERROR + "invalid role")
	}

	// set PORT2 to PORT + 10000
	reg := regexp.MustCompile(`(.*)((?::))((?:[0-9]+))$`)
	pieces := reg.FindStringSubmatch(serverConfig.Listen)
	if len(pieces) != 4 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR + serverConfig.Listen)
	}
	iPort, err := strconv.Atoi(pieces[3])
	if err != nil {
		return nil, errors.New(resp.RES_SYNTAX_ERROR + serverConfig.Listen)
	}
	serverConfig.Port2 = strconv.Itoa(10000 + iPort)
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
