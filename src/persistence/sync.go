/*
* author: duan hongxing
* date: 2 Aug, 2024
* desc:
*	Sync log data from Partition Leader
 */

package persistence

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
	"topiik/cluster"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/resp"
)

var ticker *time.Ticker
var ptnLeaderId string
var ptnLeaderAddr string
var conn *net.TCPConn

func Sync() {
	ticker = time.NewTicker(1 * time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	for {
		<-ticker.C
		if node.GetNodeInfo().Id == ptnLeaderId { // if current node is Leader, do nothing
			break
		}
		doSync()
	}

}

func doSync() {
	//l.Info().Msg("doSync")
	if len(ptnLeaderAddr) == 0 {
		err := getPartitionLeader()
		if err != nil {
			return
		}
	}
	var err error
	if conn == nil {
		conn, err = util.PreapareSocketClient(ptnLeaderAddr)
		if err != nil {
			l.Err(err).Msg(err.Error())
			ptnLeaderId = ""
			ptnLeaderAddr = ""
			return
		}
	}
	

}

/*
* Get Partition Leader from Controller
 */
func getPartitionLeader() error {
	clAddr := cluster.GetNodeStatus().LeaderControllerAddr
	if len(clAddr) == 0 {
		return errors.New("")
	}

	hostPort, err := util.SplitAddress(clAddr)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return err
	}

	conn, err := util.PreapareSocketClient(hostPort[0] + ":" + hostPort[2])
	if err != nil {
		return err
	}
	defer conn.Close()

	var req []byte
	var byteBuf = new(bytes.Buffer) // int to byte byte buf
	// 1 bytes of command
	binary.Write(byteBuf, binary.LittleEndian, cluster.RPC_GET_PL)
	req = append(req, byteBuf.Bytes()...)
	req = append(req, []byte(node.GetNodeInfo().Id)...) // send current worker node id to get leader id

	// enocde
	req, err = proto.EncodeB(req)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return err
	}

	// send
	_, err = conn.Write(req)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return err
	}

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			l.Err(err).Msgf("sync::getPartitionLeader %s", err.Error())
		}
		return err
	}
	if len(buf) < resp.RESPONSE_HEADER_SIZE {
		l.Warn().Msgf("sync::getPartitionLeader invalid response len%v", len(buf))
		return errors.New("")
	}

	// read response flag
	byteBuf.Reset()
	bufSlice := buf[4:5]
	byteBuf = bytes.NewBuffer(bufSlice)
	var flag int8
	err = binary.Read(byteBuf, binary.LittleEndian, &flag)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return err
	}
	if flag == 1 {
		idAndAddr2 := string(buf[resp.RESPONSE_HEADER_SIZE:])
		ptnLeaderId = idAndAddr2[:10]
		ptnLeaderAddr = idAndAddr2[10:]
	}
	l.Info().Msgf("sync::getPartitionLeader id, addr2: %s, %s", ptnLeaderId, ptnLeaderAddr)
	return nil
}
