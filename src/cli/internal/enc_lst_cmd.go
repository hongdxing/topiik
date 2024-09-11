//author: Duan Hongxing
//date: 11 Sep, 2024

package internal

import (
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
)

func encLslice(pieces []string) ([]byte, error) {
	if len(pieces) != 3 {
		return syntaxErr()
	}
	keys := Abytes{[]byte(pieces[1])}
	vals := Abytes{}
	args := strings.Join(pieces[2:], consts.SPACE)

	build := &CmdBuilder{Cmd: command.LSLICE_I, Ver: 1}
	return build.BuildM(keys, vals, args)
}
