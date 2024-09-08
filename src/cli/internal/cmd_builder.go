//author: Duan HongXing
//date: 25 Aug, 2024

package internal

import (
	"errors"
	"strings"
	"topiik/internal/proto"
	"topiik/resp"
)

type CmdBuilder struct {
	Ver uint8
	Cmd uint8
}

func (c *CmdBuilder) Build(key string, val string, input string) (buf []byte, err error) {
	buf, err = proto.EncodeHeader(c.Cmd, c.Ver)
	if err != nil {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	args := strings.TrimSpace(input[len(key)+len(val)+1:])
	req := &Req{}
	data, err := req.WithKey(key).WithVal(val).WithArgs(args).Marshal()
	if err != nil {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	buf = append(buf, data...)
	return buf, nil
}

func (c *CmdBuilder) BuildM(keys Abytes, vals Abytes, args string) (buf []byte, err error) {
	buf, err = proto.EncodeHeader(c.Cmd, c.Ver)
	if err != nil {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	req := &Req{}
	data, err := req.WithKeys(keys).WithVals(vals).WithArgs(args).Marshal()
	if err != nil {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	buf = append(buf, data...)
	return buf, nil
}
