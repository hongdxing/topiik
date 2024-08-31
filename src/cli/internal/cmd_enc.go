package internal

import (
	"encoding/json"
	"errors"
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/proto"
	"topiik/internal/util"
)

/*
* Encode string command to byte array
*
 */
/*
func EncodeCmd(strCmd string) ([]byte, error) {
	if icmd, ok := command.CmdCode[strCmd]; ok {
		byteBuf := new(bytes.Buffer)
		binary.Write(byteBuf, binary.LittleEndian, icmd)
		return byteBuf.Bytes(), nil
	}
	return nil, errors.New("invalid command")
}
*/

const syntax_err = "syntax error"

func EncodeCmd(input string, theCMD *string) (result []byte, err error) {
	/*pieces := strings.SplitN(input, consts.SPACE, 2)
	if len(pieces) != 2 {
		return nil, errors.New("sytax error")
	}*/
	pieces, err := util.SplitCommandLine(input)
	if err != nil {
		return nil, err
	}
	cmd := strings.ToUpper(pieces[0])
	theCMD = &cmd
	var icmd uint8 = 0
	var ver uint8 = 0
	req := datatype.Req{Keys: datatype.Abytes{}, Vals: datatype.Abytes{}}
	if cmd == command.INIT_CLUSTER { // INIT-CLUSTER partitions count
		if len(pieces) != 3 {
			return syntaxErr()
		}
		if strings.ToLower(pieces[1]) != "partition" {
			return syntaxErr()
		}
		icmd = command.INIT_CLUSTER_I
		req.Args = strings.Join(pieces[2:], consts.SPACE)
	} else if cmd == command.ADD_CONTROLLER { // ADD-CONTROLLER host:port
		if len(pieces) != 2 {
			return syntaxErr()
		}
		icmd = command.ADD_CONTROLLER_I
		req.Args = strings.Join(pieces[1:], consts.SPACE)
	} else if cmd == command.ADD_WORKER { // ADD-WORKER host:port [partition xxx]
		if len(pieces) < 3 {
			return syntaxErr()
		}
		icmd = command.ADD_WORKER_I
		req.Args = strings.Join(pieces[1:], consts.SPACE)
	} else if cmd == command.SHOW {
		return encShowCluster(pieces)
	} else if cmd == command.SCALE {
		/*
		* Syntax: scale partition 1 replica 2
		 */
		if len(pieces) != 5 {
			return syntaxErr()
		}
		icmd = command.SCALE_I
		req.Args = strings.Join(pieces[1:], consts.SPACE)
	} else if cmd == command.SET {
		return encSET(pieces)
	} else if cmd == command.GET {
		return encGET(pieces)
	} else if cmd == command.SETM { // SETM k1 v1 k2 v2
		if len(pieces) < 3 || (len(pieces)-1)%2 != 0 {
			return syntaxErr()
		}
		icmd = command.SETM_I
		for i := 1; i < len(pieces)-1; i += 2 {
			req.Keys = append(req.Keys, []byte(pieces[i]))
			req.Vals = append(req.Vals, []byte(pieces[i+1]))
		}
	} else if cmd == command.GETM { // GETM k1 k2
		if len(pieces) < 2 {
			return syntaxErr()
		}
		icmd = command.GETM_I
		for _, piece := range pieces[1:] {
			req.Keys = append(req.Keys, []byte(piece))
		}
		//req.Keys = append(req.Keys, pieces[1:]...)
	} else if cmd == command.TTL {
		if len(pieces) < 1 {
			return syntaxErr()
		}
		icmd = command.TTL_I
		req.Keys = append(req.Keys, []byte(pieces[0]))
		req.Args = strings.Join(pieces[1:], consts.SPACE)
	} else if cmd == command.LPUSH { /* LIST COMMANDS START */
		/*
		* Syntax: LPUSH key v1 [v2 v3 ...]
		 */
		if len(pieces) < 3 {
			return syntaxErr()
		}
		icmd = command.LPUSH_I
		req.Keys = append(req.Keys, []byte(pieces[1]))
		for _, piece := range pieces[2:] {
			req.Vals = append(req.Vals, []byte(piece))
		}
		//req.Vals = append(req.Vals, pieces[2:]...)
	} else if cmd == command.LPOP {
		/*
		* Syntax: LPOP key [count]
		 */
		if len(pieces) < 2 {
			return syntaxErr()
		}
		icmd = command.LPOP_I
		req.Keys = append(req.Keys, []byte(pieces[1]))
		req.Args = strings.Join(pieces[2:], consts.SPACE)
	} else if cmd == command.LPUSHR {
		/*
		* Syntax: LPUSHR key v1 [v2 v3 ...]
		 */
		if len(pieces) < 3 {
			return syntaxErr()
		}
		icmd = command.LPUSHR_I
		req.Keys = append(req.Keys, []byte(pieces[1]))
		for _, piece := range pieces[2:] {
			req.Vals = append(req.Vals, []byte(piece))
		}
		//req.Vals = append(req.Vals, pieces[2:]...)
	} else if cmd == command.LPOPR {
		/*
		* Syntax: LPOP key [count]
		 */
		if len(pieces) < 2 {
			return syntaxErr()
		}
		icmd = command.LPOPR_I
		req.Keys = append(req.Keys, []byte(pieces[1]))
		req.Args = strings.Join(pieces[2:], consts.SPACE)
	} else if cmd == command.LLEN {
		if len(pieces) != 2 {
			return syntaxErr()
		}
		icmd = command.LLEN_I
		req.Keys = append(req.Keys, []byte(pieces[1]))
	} else if cmd == command.GET_CTLADDR { /* Get Cluster Leader Addr for client to redirect */
		// no additional data
	} else if cmd == command.Keys { /* KEY COMMANDS START */
		if len(pieces) != 2 {
			return syntaxErr()
		}
		req.Args = pieces[1]
		icmd = command.KEYS_I
	} else {
		return nil, errors.New("syntax error")
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	result, err = proto.EncodeHeader(icmd, ver)
	if err != nil {
		return nil, errors.New("syntax error")
	}
	result = append(result, reqBytes...)

	return result, nil

}

/*Cluster--------------------------------------------------------------------*/
func encShowCluster(pieces []string) ([]byte, error) {
	cmdBuilder := CmdBuilder{Cmd: command.SHOW_I, Ver: 1}
	return cmdBuilder.BuildM(Abytes{}, Abytes{}, "")
}

/*String---------------------------------------------------------------------*/
func encSET(pieces []string) ([]byte, error) {
	if len(pieces) < 3 {
		return syntaxErr()
	}
	keys := Abytes{[]byte(pieces[1])}
	vals := Abytes{[]byte(pieces[2])}
	args := ""
	if len(pieces) > 3 {
		args = strings.Join(pieces[3:], consts.SPACE)
	}
	cmdBuilder := &CmdBuilder{Cmd: command.SET_I, Ver: 1}
	return cmdBuilder.BuildM(keys, vals, args)
}

func encGET(pieces []string) ([]byte, error) {
	if len(pieces) < 2 {
		return syntaxErr()
	}
	keys := Abytes{[]byte(pieces[1])}
	vals := Abytes{}
	args := ""
	if len(pieces) > 2 {
		args = strings.Join(pieces[2:], consts.SPACE)
	}
	cmdBuilder := &CmdBuilder{Cmd: command.GETM_I, Ver: 1}
	return cmdBuilder.BuildM(keys, vals, args)
}

/*List-----------------------------------------------------------------------*/

func errResult(e string) ([]byte, error) {
	return nil, errors.New(e)
}

func syntaxErr() ([]byte, error) {
	return errResult(syntax_err)
}