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
	"topiik/internal/config"
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
func clusterJoin(pieces []string, serverConfig *config.ServerConfig) (result string, err error) {

	conn, err := util.PreapareSocketClient(pieces[0])
	if err != nil {
		return "", errors.New("Join to cluster failed, please check whether captial node still alive and try again")
	}
	line := "CLUSTER " + command.CLUSTER_JOIN_ACK + " aaaaa " + pieces[0]

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
		return "", errors.New("Join to cluster failed, please check whether captial node still alive and try again")
	} else {
		fmt.Println(string(buf[0:n]))
		return RES_OK, nil
	}
}
