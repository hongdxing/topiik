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
func clusterJoin(addr string) (result string, err error) {

	fmt.Printf("clusterJoin:: %s\n", addr)
	conn, err := util.PreapareSocketClient(addr)
	if err != nil {
		return "", errors.New("join to cluster failed, please check whether captial node still alive and try again")
	}
	defer conn.Close()

	line := "CLUSTER " + command.CLUSTER_JOIN_ACK + " aaaaa " + addr

	data, err := proto.Encode(line)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("11111111")
	// Send
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("22222222")
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	fmt.Println("33333333")
	if err != nil {
		fmt.Println("44444444")
		return "", errors.New("join to cluster failed, please check whether captial node still alive and try again")
	} else {
		fmt.Println("555555555")
		fmt.Println(string(buf[0:n]))
		return "OK", nil
	}
}
