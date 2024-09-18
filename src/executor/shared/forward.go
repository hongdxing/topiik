//author: Duan HongXing
//data: 4 Jul, 2024

package shared

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"slices"
	"topiik/cluster"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/persistence"
	"topiik/resp"
)

// conn cache from leader to leader
var leaderConnCache = make(map[string]*net.TCPConn)

type ExeFn func(uint8, datatype.Req) ([]byte, error)

var persistCmds = []uint8{
	// String
	command.SET_I,
	command.SETM_I,
	command.INCR_I,
	// List
	command.LPUSH_I,
	command.LPUSHR_I,
	command.LPOP_I,
	command.LPOPR_I,

	//command.LPUSHB_I,
	//command.LPUSHRB_I,
	command.DEL_I,
	command.TTL_I, //??
}

func ExecuteOrForward(targetWorker node.NodeSlim, execute ExeFn, icmd uint8, req datatype.Req, msg []byte) (finalRes []byte) {
	if targetWorker.Id == node.GetNodeInfo().Id {
		finalRes, err := execute(icmd, req)

		if err != nil {
			// enqueue persistor queue
			if slices.Contains(persistCmds, icmd) {
				persistence.Enqueue(msg)
			}

			// sync to partition follower(s)
			// Q: what if follower down???
			// Q: what if follower fall behind???
			ptn := cluster.GetPtnByNodeId(node.GetNodeInfo().Id)
			if len(ptn.Nodes) > 0 {
				var ptnFlrs []node.NodeSlim
				for _, nd := range ptn.Nodes {
					if nd.Id != node.GetNodeInfo().Id {
						ptnFlrs = append(ptnFlrs, nd)
					}
				}
				if len(ptnFlrs) > 0 {
					// TODO: retry sync to follower
					err = persistence.SyncFollower(ptnFlrs, msg)
					if err != nil {
						return resp.ErrResponse(fmt.Errorf("%s %s", "todo", "todo"))
					}
				}
			}
		}

		return finalRes
	} else {
		//return shared.ForwardByKey(key, msg, targetWorker)
		return ForwardByWorker(targetWorker, msg)
	}
}

func ForwardByKey(key []byte, msg []byte, targetWorker node.NodeSlim) []byte {
	var err error
	if len(targetWorker.Id) == 0 {
		l.Warn().Msg("forward::Forward no slot available")
		return resp.ErrResponse(errors.New(resp.RES_NO_WORKER))
	}

	conn, ok := leaderConnCache[targetWorker.Id]
	if !ok {
		//conn, err = util.PreapareSocketClientWithPort(targetWorker.Addr, CONTROLLER_FORWORD_PORT)
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
		leaderConnCache[targetWorker.Id] = conn
	}
	// Send
	_, err = conn.Write(msg)
	if err != nil {
		l.Err(err).Msg(err.Error())
		if _, ok = leaderConnCache[targetWorker.Id]; ok {
			l.Warn().Msgf("forward::Forward remove tcp cache of worker %s", targetWorker.Id)
			delete(leaderConnCache, targetWorker.Id)
		}
		// try reconnect
		targetWorker, err := GetLeaderNode(key) // the worker may changed because of Worker Leader fail
		if err != nil {
			return resp.ErrResponse(err)
		}
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
		conn.Write(msg)
	}

	reader := bufio.NewReader(conn)
	responseBytes, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
	}
	return responseBytes
}

func GetLeaderNode(key []byte) (node node.NodeSlim, err error) {
	var keyHash = crc32.Checksum(key, crc32.IEEETable)
	keyHash = keyHash % consts.SLOTS
	return cluster.GetNodeByKeyHash(uint16(keyHash))
}

func ForwardByWorker(targetWorker node.NodeSlim, msg []byte) []byte {
	var err error
	conn, ok := leaderConnCache[targetWorker.Id]
	if !ok {
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
		leaderConnCache[targetWorker.Id] = conn
	}
	// Send
	_, err = conn.Write(msg)
	if err != nil {
		l.Err(err).Msg(err.Error())
		if _, ok = leaderConnCache[targetWorker.Id]; ok {
			l.Warn().Msgf("forward::Forward remove tcp cache of worker %s", targetWorker.Id)
			delete(leaderConnCache, targetWorker.Id)
		}
		// try reconnect
		//targetWorker := getWorker(key) // the worker may changed because of Worker Leader fail
		conn, err = util.PreapareSocketClient(targetWorker.Addr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
		conn.Write(msg)
	}
	reader := bufio.NewReader(conn)
	responseBytes, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			return resp.ErrResponse(errors.New(resp.RES_CONN_RESET))
		}
	}
	return responseBytes
}

func syncFollower() {

}
