package persistence

import (
	"bytes"
	"encoding/binary"
	"io"
)

/*
* Read one log from binary log
*
 */
func parseOne(r io.Reader, seq *int64) (res []byte, err error) {
	//var seq int64
	var length int32

	// 1: read seq
	buf := make([]byte, 8)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}

	byteBuf := bytes.NewBuffer(buf)
	err = binary.Read(byteBuf, binary.LittleEndian, seq)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return
	}
	res = append(res, buf...)

	// 2: read length
	buf = make([]byte, 4)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	byteBuf = bytes.NewBuffer(buf)
	binary.Read(byteBuf, binary.LittleEndian, &length)
	res = append(res, buf...)

	// 3: read msg
	buf = make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	res = append(res, buf...)
	return res, nil
}
