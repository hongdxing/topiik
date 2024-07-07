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
	"regexp"
	"strings"
	"topiik/cluster"
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
**	CLUSTER JOIN host:port CONTROLLER|WORKER
**/
func clusterJoin(myAddr string, controllerAddr string, role string) (result string, err error) {
	reg, _ := regexp.Compile(consts.EMAIL_PATTERN)
	if !reg.MatchString(myAddr) || !reg.MatchString(controllerAddr) {
		println("xxxx")
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	if strings.ToUpper(role) != cluster.ROLE_CONTROLLER && strings.ToUpper(role) != cluster.ROLE_WORKER {
		return "", errors.New(RES_SYNTAX_ERROR)
	}
	nodeId := cluster.GetNodeMetadata().Id

	fmt.Printf("clusterJoin:: %s\n", controllerAddr)
	conn, err := util.PreapareSocketClient(controllerAddr)
	if err != nil {
		return "", errors.New("join to cluster failed, please check whether captial node still alive and try again")
	}
	defer conn.Close()

	// CLUSTER JOIN_ACK nodeId addr role
	line := "CLUSTER " + command.CLUSTER_JOIN_ACK + consts.SPACE + nodeId + consts.SPACE + myAddr + consts.SPACE + role
	fmt.Println(line)

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
		fmt.Println(string(buf[:n]))
		err = cluster.UpdateNodeClusterId(string(buf[1:n])) // the first byte is flag
		if err != nil {
			fmt.Println(err)
			return "", errors.New("join cluster failed")
		}
		// if join controller succeed, will start to RequestVote
		if strings.ToUpper(role) == cluster.ROLE_CONTROLLER {
			go cluster.RequestVote()
		}
		return RES_OK, nil
	}
}
