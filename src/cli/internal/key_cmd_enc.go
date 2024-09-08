//author: Duan HongXing
//date: 8 Sep, 2024

package internal

import "topiik/internal/command"

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
	return builder.BuildM(keys, Abytes{}, pieces[1])
}
