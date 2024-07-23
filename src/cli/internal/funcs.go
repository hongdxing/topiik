package internal

import (
	"encoding/json"
	"errors"
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/internal/util"
)

/*
** Encode string command to byte array
**
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

func EncodeCmd(input string) (result []byte, err error) {
	/*pieces := strings.SplitN(input, consts.SPACE, 2)
	if len(pieces) != 2 {
		return nil, errors.New("sytax error")
	}*/
	pieces, err := util.SplitCommandLine(input)
	if err != nil {
		return nil, err
	}
	cmd := strings.ToUpper(pieces[0])
	req := datatype.Req{VER: 1, CMD: cmd, KEYS: []string{}, VALS: []string{}}
	if cmd == command.INIT_CLUSTER { // INIT-CLUSTER 1 1
		if len(pieces) != 3 {
			return syntaxErr()
		}
		req.ARGS = strings.Join(pieces[1:], consts.SPACE)
	} else if cmd == command.ADD_NODE { // ADD-NODE host:port role
		if len(pieces) != 3 {
			return syntaxErr()
		}
		req.ARGS = strings.Join(pieces[1:], consts.SPACE)
	} else if cmd == command.SET { // SET key val args
		if len(pieces) < 3 {
			return syntaxErr()
		}
		req.KEYS = []string{pieces[1]} // key
		req.VALS = []string{pieces[2]} // value
		// the rest as args
		if len(pieces) > 3 {
			req.ARGS = strings.Join(pieces[3:], consts.SPACE)
		}
	} else if cmd == command.GET { // GET k1
		if len(pieces) != 2 {
			return syntaxErr()
		}
		req.KEYS = []string{pieces[1]} // key
	} else if cmd == command.SETM { // SETM k1 v1 k2 v2
		if len(pieces) < 3 || (len(pieces)-1)%2 != 0 {
			return syntaxErr()
		}
		for i := 1; i < len(pieces)-1; i += 2 {
			req.KEYS = append(req.KEYS, pieces[i])
			req.VALS = append(req.VALS, pieces[i+1])
		}
	} else if cmd == command.GETM { // GETM k1 k2
		if len(pieces) < 2 {
			return syntaxErr()
		}
		req.KEYS = append(req.KEYS, pieces[1:]...)
	} else if cmd == command.GET_CTL_LEADER_ADDR {
		// no additional data
	} else {
		return nil, errors.New("syntax error")
	}
	result, err = json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func errResult(e string) ([]byte, error) {
	return nil, errors.New(e)
}

func syntaxErr() ([]byte, error) {
	return errResult(syntax_err)
}
