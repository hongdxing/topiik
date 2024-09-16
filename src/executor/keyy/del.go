//author: Duan HongXing
//date: 27 Aug, 2024

package keyy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"topiik/cluster"
	"topiik/executor/shared"
	"topiik/internal/command"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

// Delete the key(s) if exists
// Return number of key(s) deleted
func Del(req datatype.Req) (rslt int64, err error) {
	for _, keyB := range req.Keys {
		key := string(keyB)
		if _, ok := memo.MemMap[string(key)]; ok {
			delete(memo.MemMap, key)
			rslt++
		}
	}
	return rslt, err
}

// forward DEL command to each group
func ForwardDel(execute shared.ExeFn, req datatype.Req, msg []byte) (rslt int64) {
	var err error
	for _, worker := range cluster.GetPtnLeaders() {
		//buf := shared.ForwardByWorker(worker, msg) // get keys from each worker leader
		buf := shared.ExecuteOrForward(worker, execute, command.DEL_I, req, msg)
		if len(buf) > 4 {
			bbuf := bytes.NewBuffer(buf[4:5])
			var flag resp.RespFlag
			err = binary.Read(bbuf, binary.LittleEndian, &flag)
			if err != nil {
				l.Warn().Msgf("keys::ForwardDel %s", err.Error())
				continue
			}

			if flag == resp.Success {
				bbuf = bytes.NewBuffer(buf[5:6])
				var datatype resp.RespType
				err = binary.Read(bbuf, binary.LittleEndian, &datatype)
				if err != nil {
					l.Warn().Msgf("keys::ForwardDel %s", err.Error())
					continue
				}

				if datatype == resp.Integer {
					var partialRes int64
					bbuf = bytes.NewBuffer(buf[resp.RESPONSE_HEADER_SIZE:])
					binary.Read(bbuf, binary.LittleEndian, &partialRes)

					if err != nil {
						fmt.Printf("(err):%s\n", err.Error())
					}
					rslt += partialRes
				} else {
					fmt.Println("(err): invalid response type")
				}

			} else {
				res := buf[resp.RESPONSE_HEADER_SIZE:]
				l.Warn().Msgf("keys::ForwardDel %s", string(res))
			}
		}
	}
	return rslt
}
