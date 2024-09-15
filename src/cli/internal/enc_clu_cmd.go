//author: Duan HongXing
//date: 9 Sep, 2024

package internal

import (
	"strings"
	"topiik/internal/command"
	"topiik/internal/consts"
)

func encInitCluster(pieces []string) ([]byte, error) {
	builder := CmdBuilder{Cmd: command.CREATE_CLUSTER_I, Ver: 1}
	// The first piece must be partition
	if strings.ToLower(pieces[1]) != "partition" {
		return syntaxErr()
	}
	args := strings.Join(pieces[1:], consts.SPACE)
	return builder.BuildM(Abytes{}, Abytes{}, args)
}

func encShowCluster(pieces []string) ([]byte, error) {
	cmdBuilder := CmdBuilder{Cmd: command.SHOW_I, Ver: 1}
	return cmdBuilder.BuildM(Abytes{}, Abytes{}, "")
}

func encRemoveNode(pieces []string) ([]byte, error) {
	if len(pieces) != 2 {
		return syntaxErr()
	}
	builder := CmdBuilder{Cmd: command.REMOVE_NODE_I, Ver: 1}
	return builder.BuildM(Abytes{}, Abytes{}, pieces[1])
}

func encNewPartition(pieces []string) ([]byte, error) {
	builder := CmdBuilder{Cmd: command.NEW_PARTITION_I, Ver: 1}
	return builder.BuildM(Abytes{}, Abytes{}, "")
}
