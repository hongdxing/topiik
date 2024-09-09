//author: Duan HongXing
//date: 9 Sep, 2024

package internal

import (
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
)

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
	cmdBuilder := &CmdBuilder{Cmd: command.GET_I, Ver: 1}
	return cmdBuilder.BuildM(keys, vals, args)
}
