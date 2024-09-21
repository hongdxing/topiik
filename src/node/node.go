package node

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"topiik/internal/config"
	"topiik/internal/consts"
	"topiik/internal/util"
	"topiik/resp"
)

const (
	slash   = string(os.PathSeparator)
	dataDIR = "data"
)

const cluster_init_failed = "cluster init failed"

func InitNode(serverConfig config.ServerConfig) (err error) {
	l.Info().Msg("node::InitNode start")

	var exist bool
	fpath := GetNodeFilePath()

	// data dir
	exist, err = util.PathExists(dataDIR)
	if err != nil {
		l.Panic().Msg(err.Error())
	}

	if !exist {
		l.Info().Msg("node::InitNode Creating data dir...")
		err = os.Mkdir(dataDIR, os.ModeDir)
		if err != nil {
			l.Panic().Msg(err.Error())
		}
	}

	// node file
	exist, err = util.PathExists(fpath)
	if err != nil {
		l.Panic().Msg(err.Error())
	}

	var buf []byte
	var node Node
	if !exist {
		l.Info().Msg("node::InitNode creating node file...")

		node.Id = strings.ToLower(util.RandStringRunes(consts.NODE_ID_LEN))
		node.Role = strings.ToLower(serverConfig.Role)
		node.Addr = serverConfig.Listen
		node.Addr2 = serverConfig.Host + ":" + serverConfig.Port2
		buf, _ = json.Marshal(node)
		err = os.WriteFile(fpath, buf, 0644)
		if err != nil {
			l.Panic().Msg("node::InitNode loading node failed")
		}
	} else {
		l.Info().Msg("node::InitNode loading node...")

		buf, err = os.ReadFile(fpath)
		if err != nil {
			l.Panic().Msg("node::InitNode loading node failed")
		}
		err = json.Unmarshal(buf, &node)
		if err != nil {
			l.Panic().Msg("node::InitNode loading node failed")
		}
		node.Addr = serverConfig.Listen
		node.Addr2 = serverConfig.Host + ":" + serverConfig.Port2
		l.Info().Msgf("node::InitNode load node %s", node)
	}
	nodeInfo = &node

	l.Info().Msgf("node::InitNode node role: %s", nodeInfo.Role)
	l.Info().Msg("node::InitNode end")
	return nil
}

// Update Node clusterId when init cluster
func InitCluster(clusterId string, ptnId string) (err error) {
	// 1. open node file
	fpath := GetNodeFilePath()

	// 2. read from node file
	var jsonStr string
	var buf []byte

	/*buf, err = os.ReadFile(fpath)
	if err != nil {
		return errors.New(cluster_init_failed)
	}*/
	buf, err = util.ReadBinaryFile(fpath)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return errors.New(resp.RES_OPEN_FILE_ERR)
	}
	jsonStr = string(buf)

	// 3. unmarshal
	err = json.Unmarshal([]byte(jsonStr), &nodeInfo)
	if err != nil {
		return errors.New(cluster_init_failed)
	}
	if len(nodeInfo.ClusterId) > 0 { // check if current node already in cluster or not
		return errors.New("current node already in a cluster:" + nodeInfo.ClusterId)
	}
	nodeInfo.ClusterId = clusterId
	nodeInfo.PntId = ptnId
	nodeInfo.Role = config.ROLE_WORKER

	// 4. marshal
	buf2, _ := json.Marshal(nodeInfo)

	// 5. write back to file
	util.WriteBinaryFile(fpath, buf2)
	//os.Truncate(fpath, 0)
	//os.WriteFile(fpath, buf2, 0644)

	return nil
}

// Initialized by command ADD-WORKER|ADD-CONTROLLER
// Update clusterId when Controller try to add current node to the cluster
func JoinCluster(clusterId string, ptnId string) (err error) {
	/* update node cluster id and role */
	nodeInfo.ClusterId = clusterId
	nodeInfo.PntId = ptnId

	nodePath := GetNodeFilePath()
	buf, err := json.Marshal(nodeInfo)
	if err != nil {
		return errors.New("update node failed")
	}
	err = os.Truncate(nodePath, 0) // TODO: myabe need backup first
	if err != nil {
		return errors.New("update node failed")
	}
	err = os.WriteFile(nodePath, buf, 0664) // save back controller file
	if err != nil {
		return errors.New("update node failed")
	}
	return nil
}

func SetConfig(svrConfig config.ServerConfig) {
	serverConfig = &svrConfig
}

func GetConfig() config.ServerConfig {
	return *serverConfig
}

/* sync to avoid fatal error: concurrent map writes */
var setPtnMu sync.Mutex

func SetPtn(buf []byte) {
	setPtnMu.Lock()
	defer setPtnMu.Unlock()
	err := json.Unmarshal(buf, &partition)
	if err != nil {
		l.Err(err).Msg(err.Error())
	}
	//l.Info().Msg(string(buf))
}

func GetNodeInfo() Node {
	return *nodeInfo
}

func GetNodeFilePath() string {
	return util.GetMainPath() + slash + dataDIR + slash + "__metadata_node__"
}

func IsWorker() bool {
	return nodeInfo.Role == config.ROLE_WORKER
}

func IsPersistor() bool {
	return nodeInfo.Role == config.ROLE_PERSISTOR
}
