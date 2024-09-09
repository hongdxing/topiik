//author: Duan HongXing
//date: 8 Sep, 2024

package internal

import (
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
)

// DEL Key(s)
func encDEL(pieces []string) ([]byte, error) {
	// cmd + at least 1 key
	if len(pieces) < 2 {
		return syntaxErr()
	}
	builder := CmdBuilder{Cmd: command.DEL_I, Ver: 1}
	var keys Abytes
	for _, piece := range pieces[1:] {
		keys = append(keys, []byte(piece))
	}
	return builder.BuildM(keys, Abytes{}, "")
}

// Exists Key(s)
func encExists(pieces []string) ([]byte, error) {
	// cmd + at least 1 key
	if len(pieces) < 2 {
		return syntaxErr()
	}
	builder := CmdBuilder{Cmd: command.EXISTS_I, Ver: 1}
	var keys Abytes
	for _, piece := range pieces[1:] {
		keys = append(keys, []byte(piece))
	}

	return builder.BuildM(keys, Abytes{}, "")
}

func encTtl(pieces []string) ([]byte, error) {
	if len(pieces) < 2 {
		return syntaxErr()
	}
	keys := Abytes{[]byte(pieces[1])}
	args := strings.Join(pieces[2:], consts.SPACE)
	builder := CmdBuilder{Cmd: command.TTL_I, Ver: 1}
	return builder.BuildM(keys, Abytes{}, args)
}
