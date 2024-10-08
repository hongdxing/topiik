//author: Duan HongXing
//date: 19 Jun, 2024

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"topiik/cli/internal"
	"topiik/internal/command"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/resp"
)

var host string
var pass string

// Syntax: topiik-cli(.exe) --host localhost:8301 [--pass password]
func main() {
	var (
		cluster    string
		workers    string
		partitions uint
		host       string
	)

	flag.StringVar(&cluster, "cluster", "", "topiik-cli --cluster create")
	flag.StringVar(&workers, "workers", "", "")
	flag.UintVar(&partitions, "partitions", 1, "")
	flag.StringVar(&host, "host", "localhost:8301", "")
	flag.Parse()

	if strings.TrimSpace(cluster) != "" {
		cluster = strings.ToLower(strings.TrimSpace(cluster))
		if cluster == "create" {
			if strings.TrimSpace(workers) == "" {
				fmt.Println("missing workers")
			}
		}
		err := internal.CreateCluster(workers, partitions)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	// Connect to the server
	pass = ""
	const invalidArgs = "invalid args"
	fmt.Println(os.Args)
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if strings.ToLower(os.Args[i]) == "--host" {
				if len(os.Args) < i+1 {
					fmt.Println(invalidArgs)
					return
				}
				host = os.Args[i+1]
				i++
			} else if strings.ToLower(os.Args[i]) == "--pass" {
				if len(os.Args) < i+1 {
					fmt.Println(invalidArgs)
					return
				}
				pass = os.Args[i+1]
				i++
				fmt.Println(pass)
			} else {
				fmt.Println(invalidArgs)
			}
		}
	}

	// Get Controller Leader address
	leaderAddr, err := getControllerLeaderAddr(host)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// connect to leader
	conn, err := util.PreapareSocketClient(leaderAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	if leaderAddr != host {
		fmt.Printf("redirect to %s\n", leaderAddr)
	}

	reader := bufio.NewReader(os.Stdin)

	// for loop keep cli alive
	for {
		fmt.Print(leaderAddr + ">")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, " \t\r\n")
		if len(line) == 0 {
			continue
		}

		/*
			pieces := strings.SplitN(line, " ", 2)
			cmd := strings.ToUpper(pieces[0])
			if cmd == command.S_QUIT {
				conn.Close()
				break
			}
			msg, err := internal.EncodeCmd(cmd)
			if err != nil {
				fmt.Println(err.Error())
			}

			if len(pieces) == 2 {
				msg = append(msg, []byte(pieces[1])...)
			}
		*/

		// TODO: valid command
		// Enocde
		var cmd string
		msg, err := internal.EncodeCmd(line, &cmd)
		data, err := proto.EncodeB(msg)
		if err != nil {
			fmt.Println()
			fmt.Println(err)
		}

		// Send
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println()
			fmt.Println(err)
			return
		}

		err = response(conn, cmd)
		if err != nil {
			break
		}
	}
}

func response(conn *net.TCPConn, cmd string) error {
	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("(err): %s\n", err)
		}
		return err
	}

	if len(buf) > 4 {
		bufSlice := buf[4:5]
		bbuf := bytes.NewBuffer(bufSlice)
		var flag resp.RespFlag
		err = binary.Read(bbuf, binary.LittleEndian, &flag)
		if err != nil {
			fmt.Println("(err):")
		}

		if flag == resp.Success {
			bufSlice = buf[5:6]
			bbuf = bytes.NewBuffer(bufSlice)
			var resType resp.RespType
			err = binary.Read(bbuf, binary.LittleEndian, &resType)
			if err != nil {
				fmt.Println("(err):")
			}
			if resType == resp.String {
				val := buf[resp.RESPONSE_HEADER_SIZE:]
				if string(val) == proto.Nil {
					fmt.Printf("(err):nil\n")
				} else {
					fmt.Printf("%s\n", val)
				}
			} else if resType == resp.StringArray {
				bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
				var vals []string

				err = json.Unmarshal(bufSlice, &vals)
				if err != nil {
					fmt.Printf("(err):%s\n", err.Error())
				}
				if cmd == command.CREATE_CLUSTER {
					fmt.Println("Partitions:")
				}
				if len(vals) > 0 {
					for i, v := range vals {
						if string(v) == proto.Nil {
							fmt.Printf("%v) (err):nil\n", i+1)
						} else {
							fmt.Printf("%v) %s\n", i+1, v)
						}
					}
				} else {
					fmt.Println("[]")
				}
			} else if resType == resp.Integer {
				bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
				bbuf := bytes.NewBuffer(bufSlice)
				var result int64
				err = binary.Read(bbuf, binary.LittleEndian, &result)
				if err != nil {
					fmt.Println("(err):")
				}
				fmt.Printf("%v\n", result)
			} else if resType == resp.Map {
				//
			} else {
				fmt.Println("(err): invalid response type")
			}
			/*else if resType == resp.ByteArray {
				bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
				var result datatype.Abytes

				err = json.Unmarshal(bufSlice, &result)
				if err != nil {
					fmt.Printf("(err):%s\n", err.Error())
				}
				for i, v := range result {
					fmt.Printf("%v: %s\n", i, string(v))
				}
				//fmt.Printf("%v\n", result)
			} */
		} else {
			res := buf[resp.RESPONSE_HEADER_SIZE:]
			fmt.Printf("(err):%s\n", res)
		}
	} else {
		fmt.Println("(err): unknown")
		return errors.New("unknown")
	}
	return nil
}

// get controller leader address
// return:
//   - leader address
func getControllerLeaderAddr(host string) (addr string, err error) {
	conn, err := util.PreapareSocketClient(host)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	/*
		byteBuf := new(bytes.Buffer)
		binary.Write(byteBuf, binary.LittleEndian, command.GET_CONTROLLER_LEADER_ADDR)
		data, _ := proto.EncodeB(byteBuf.Bytes())
		conn.Write(data)
	*/
	/*
		msg, err := internal.EncodeCmd(command.GET_LEADER_ADDR)
		if err != nil {
			fmt.Printf("(err):%s\n", err)
		}*/
	msg, err := proto.EncodeHeader(command.GET_CTLADDR_I, 1)
	if err != nil {
		return "", errors.New("syntax error")
	}
	req := datatype.Req{Keys: datatype.Abytes{}, Vals: datatype.Abytes{}}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("error")
	}
	msg = append(msg, reqBytes...)
	data, _ := proto.EncodeB(msg)
	conn.Write(data)

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("(err): %s\n", err)
		}
		return "", err
	}
	//fmt.Println(string(buf[6:]))
	if len(buf) > 4 {
		bufSlice := buf[4:5]
		byteBuf := bytes.NewBuffer(bufSlice)
		var flag int8
		err = binary.Read(byteBuf, binary.LittleEndian, &flag)
		if err != nil {
			return "", err
		}

		if flag == 1 {
			bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
			return string(bufSlice), nil
		} else {
			bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
			return "", errors.New("(err):" + string(bufSlice))
		}
	}
	return "", errors.New("(err): unknown")
}
