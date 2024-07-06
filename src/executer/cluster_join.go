/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package executer

import (
	"errors"
	"fmt"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
)

/***
** Join to CCSS cluster
**
**
** Syntax:
**	CLUSTER JOIN host:port
**/
func clusterJoin(nodeId string, myAddr string, captialAddr string) (result string, err error) {

	fmt.Printf("clusterJoin:: %s\n", captialAddr)
	conn, err := util.PreapareSocketClient(captialAddr)
	if err != nil {
		return "", errors.New("join to cluster failed, please check whether captial node still alive and try again")
	}
	defer conn.Close()

	line := "CLUSTER " + command.CLUSTER_JOIN_ACK + consts.SPACE + nodeId + consts.SPACE + myAddr

	data, err := proto.Encode(line)
	if err != nil {
		fmt.Println(err)
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
	}

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return "", errors.New("join to cluster failed, please check whether captial node still alive and try again")
	} else {
		fmt.Println(string(buf[0:n]))
		return "OK", nil
	}
}
