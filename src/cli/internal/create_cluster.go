// author: Duan Hongxing
// date: 14 Sep, 2024

package internal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/resp"
)

func CreateCluster(workers string, partitions uint) error {
	addrs := strings.Split(workers, ",")
	if len(addrs)%int(partitions) != 0 {
		return errors.New(resp.RES_WRONG_ARG + " workers count and partition number not even")
	}
	_, _, err := checkConnection(addrs)
	if err != nil {
		return err
	}

	args := workers + consts.SPACE + strconv.Itoa(int(partitions))
	fmt.Printf("args: %s\n", args)
	build := &CmdBuilder{Cmd: command.CREATE_CLUSTER_I, Ver: 1}
	data, _ := build.BuildM(Abytes{}, Abytes{}, args)
	return doCreateCluster(addrs[0], data)
}

// Check connectivity of nodes
// Return id->addr and id->addr maps
func checkConnection(addrs []string) (map[string]string, map[string]string, error) {
	var addrMap = make(map[string]string)
	var addr2Map = make(map[string]string)
	for _, addr := range addrs {
		host, _, port2, err := util.SplitAddress2(strings.TrimSpace(addr))
		if err != nil {
			return nil, nil, err
		}

		addr2 := host + ":" + port2
		conn, err := util.PreapareSocketClient(addr2)
		if err != nil {
			return nil, nil, err
		}
		defer conn.Close()

		// Prepare buf
		var buf []byte
		bbuf := new(bytes.Buffer)
		binary.Write(bbuf, binary.LittleEndian, consts.RPC_TEST_CONN)
		buf = append(buf, bbuf.Bytes()...)
		buf, err = proto.EncodeB(buf)
		if err != nil {
			return nil, nil, err
		}
		// Write
		_, err = conn.Write(buf)
		if err != nil {
			return nil, nil, err

		}
		// Read
		reader := bufio.NewReader(conn)
		res, err := proto.Decode(reader)
		if err != nil {
			return nil, nil, err
		}

		// Flag
		flag := resp.ParseResFlag(res)

		if flag == resp.Success {
			if len(res) > resp.RESPONSE_HEADER_SIZE {
				ndId := string(res[resp.RESPONSE_HEADER_SIZE:])
				if err != nil {
					return nil, nil, err
				}
				addrMap[ndId] = addr
				addr2Map[ndId] = addr2
			} else {
				return nil, nil, err
			}
		} else {
			break
		}
	}

	return addrMap, addr2Map, nil
}

func doCreateCluster(addr string, data []byte) error {
	conn, err := util.PreapareSocketClient(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// encode
	data, err = proto.EncodeB(data)
	if err != nil {
		return err
	}

	// write
	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	// read
	reader := bufio.NewReader(conn)
	res, err := proto.Decode(reader)
	if err != nil {
		return err
	}

	// flag
	flag := resp.ParseResFlag(res)

	if flag == resp.Success {
		if len(res) > resp.RESPONSE_HEADER_SIZE {
			rslt := string(res[resp.RESPONSE_HEADER_SIZE:])
			fmt.Println(rslt)
		} else {
			return err
		}
	}
	return nil
}
