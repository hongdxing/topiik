/*
* author: duan hongxing
* date: 2 Aug, 2024
* desc:
*	Sync log data from Partition Leader
 */

package persistence

import (
	"time"
)

/*--------------Pull (obsoleted)-------------------*/
/*
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
	if node.GetNodeInfo().Id == ptnLeaderId { // if current node is Leader, do nothing
		return
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
	var req []byte
	var byteBuf = new(bytes.Buffer) // int to byte byte buf
	// 1 bytes of command
	binary.Write(byteBuf, binary.LittleEndian, consts.RPC_SYNC_BINLOG)
	req = append(req, byteBuf.Bytes()...)
	req = append(req, []byte(node.GetNodeInfo().Id)...)

	byteBuf.Reset()
	binary.Write(byteBuf, binary.LittleEndian, int64(2))
	req = append(req, byteBuf.Bytes()...)

	// enocde
	req, err = proto.EncodeB(req)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}

	// send
	_, err = conn.Write(req)
	if err != nil {
		l.Err(err).Msg(err.Error())
		conn = nil // clear the broken conn
		return
	}

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			l.Err(err).Msgf("persistence::getPartitionLeader %s", err.Error())
		}
		return
	}
	if len(buf) < resp.RESPONSE_HEADER_SIZE {
		l.Warn().Msgf("persistence::getPartitionLeader invalid response len%v", len(buf))
		return
	}

	// read response flag
	byteBuf.Reset()
	bufSlice := buf[4:5]
	byteBuf = bytes.NewBuffer(bufSlice)
	var flag int8
	err = binary.Read(byteBuf, binary.LittleEndian, &flag)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}
	if flag == 1 {
		buf := buf[resp.RESPONSE_HEADER_SIZE:]
	}
}
*/
/*
* Get Partition Leader from Controller
 */

/*
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
	binary.Write(byteBuf, binary.LittleEndian, consts.RPC_GET_PL)
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
}*/
/*---------------------End of Pull-----------------------*/

var ticker *time.Ticker
var syncCh = make(chan []byte)

/*
* Sync (push) to follower
*
*
 */
func Sync() {
	ticker = time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			break
			//case buf := <- syncCh:
			//	break
		}
	}
}
