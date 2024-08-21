package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
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

/***
**
**
** Syntax: topiik-cli(.exe) --host localhost:8301 [--pass password]
**/
func main() {
	// Connect to the server
	//conn, err := net.Dial("tcp", "localhost:8302")
	host = "localhost:8301"
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
		msg, err := internal.EncodeCmd(line)
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

		err = response(conn)
		if err != nil {
			break
		}
	}
}

func response(conn *net.TCPConn) error {
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
		byteBuf := bytes.NewBuffer(bufSlice)
		var flag resp.RespFlag
		err = binary.Read(byteBuf, binary.LittleEndian, &flag)
		if err != nil {
			fmt.Println("(err):")
		}

		if flag == resp.Success {
			bufSlice = buf[5:6]
			byteBuf = bytes.NewBuffer(bufSlice)
			var datatype resp.RespType
			err = binary.Read(byteBuf, binary.LittleEndian, &datatype)
			if err != nil {
				fmt.Println("(err):")
			}
			if datatype == resp.String {
				res := buf[resp.RESPONSE_HEADER_SIZE:]
				fmt.Printf("%s\n", string(res))
			} else if datatype == resp.Integer {
				bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
				byteBuf := bytes.NewBuffer(bufSlice)
				var result int64
				err = binary.Read(byteBuf, binary.LittleEndian, &result)
				if err != nil {
					fmt.Println("(err):")
				}
				fmt.Printf("%v\n", result)
			} else if datatype == resp.StringArray {
				bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
				var result []string

				err = json.Unmarshal(bufSlice, &result)
				if err != nil {
					fmt.Printf("(err):%s\n", err.Error())
				}
				fmt.Printf("%v\n", result)
			} else {
				fmt.Println("(err): invalid response type")
			}

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

/*
* Desc: get controller leader address
* Return:
*	- leader address
*
 */
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
	req := datatype.Req{KEYS: []string{}, VALS: []string{}}
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
