package persistence

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"topiik/internal/datatype"
	"topiik/internal/proto"
)

/*
* Get binlog seq
*
 */
func GetBLSeq() int64 {
	return binlogSeq
}

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

/*
* Replay binlog to load data to RAM
*
 */
func replay(buf []byte, f execute1) error {
	/* replay msg(load from disk to memory) */
	buf = buf[preMsgLen:]
	icmd, _, err := proto.DecodeHeader(buf)
	if err != nil {
		l.Err(err).Msgf("persistence::replay %s", err.Error())
		return err
	}

	var req datatype.Req
	buf = buf[2:]
	err = json.Unmarshal(buf, &req) // 2= 1 icmd and 1 ver
	if err != nil {
		l.Err(err).Msgf("persistence::replay %s", err.Error())
		return err
	}
	f(icmd, req)
	return nil
}
