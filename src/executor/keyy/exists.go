//author: Duan HongXing
//date: 28 Aug, 2024

package keyy

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"topiik/cluster"
	"topiik/executor/shared"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

// Chech if key(s) exists
// Return array of T if exist or F if not exist
func Exists(req datatype.Req) (rslt []string, err error) {
	fmt.Printf("Exists Keys %s", req.Keys)
	for _, keyB := range req.Keys {
		key := string(keyB)
		if _, ok := memo.MemMap[string(key)]; ok {
			rslt = append(rslt, "T")
		} else {
			rslt = append(rslt, "F")
		}
	}
	return rslt, err
}

// forward exists command to each partition
func ForwardExists(msg []byte, keyCount int) (rslt []string) {
	var err error
	var assemble [][]string
	rslt = make([]string, keyCount)

	// set initial value
	for i := 0; i < keyCount; i++ {
		rslt[i] = "F"
	}

	for _, worker := range cluster.GetWorkerLeaders() {
		buf := shared.ForwardByWorker(worker, msg) // get keys from each worker leader
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

				if datatype == resp.StringArray {
					var partialRes []string
					err = json.Unmarshal(buf[resp.RESPONSE_HEADER_SIZE:], &partialRes)
					if err != nil {
						l.Err(err).Msgf("(err):%s\n", err.Error())
					}
					assemble = append(assemble, partialRes)
				} else {
					l.Warn().Msg("(err): invalid response type")
				}
			} else {
				res := buf[resp.RESPONSE_HEADER_SIZE:]
				l.Warn().Msgf("keys::ForwardDel %s", string(res))
			}
		}
	}
	if len(assemble) == 1 {
		return assemble[0]
	}

	for _, arr := range assemble {
		for i, ele := range arr {
			if ele == "T" && i < keyCount {
				rslt[i] = "T"
			}
		}
	}
	return rslt
}
