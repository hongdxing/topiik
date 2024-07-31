/*
* author: duan hongxing
* date: 29 Jun 2024
* desc:
*	return keys
*
 */
package executor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"topiik/cluster"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/*
* Return keys
* Parameters:
*	- args: the arguments, command line that CMD(KEYS) stripped
* Return:
*	-
* Synctax: KEYS pattern
*	- pattern is a string to search keys, use astrisk(*) for pattern search
 */
func keys(req datatype.Req) (result []string, err error) {
	if len(req.ARGS) == 0 {
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	pattern := req.ARGS
	if !strings.HasPrefix(pattern, "*") { // exactly match from beginning
		pattern = "^" + pattern
	}
	if !strings.HasSuffix(pattern, "*") { // exactly match from endding
		pattern = pattern + "$"
	}
	//fmt.Println(strings.ReplaceAll(pattern, "*", ".*"))
	reg, err := regexp.Compile(strings.ReplaceAll(pattern, "*", ".*"))
	if err != nil {
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	keys := make([]string, 0, len(memo.MemMap))
	for k := range memo.MemMap {
		// Need to exclude internal using KEYs
		if reg.MatchString(k) && !strings.HasPrefix(k, consts.RESEVERD_PREFIX) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}

func forwardKeys(msg []byte) (res []string) {
	var err error
	for _, worker := range cluster.GetWorkerLeaders() {
		buf := cluster.ForwardByWorker(worker, msg) // get keys from each worker leader
		if len(buf) > 4 {
			bufSlice := buf[4:5]
			byteBuf := bytes.NewBuffer(bufSlice)
			var flag int8
			err = binary.Read(byteBuf, binary.LittleEndian, &flag)
			if err != nil {
				l.Warn().Msgf("keys::forwardKeys %s", err.Error())
				continue
			}

			if flag == 1 {
				bufSlice = buf[5:6]
				byteBuf = bytes.NewBuffer(bufSlice)
				var datatype int8
				err = binary.Read(byteBuf, binary.LittleEndian, &datatype)
				if err != nil {
					l.Warn().Msgf("keys::forwardKeys %s", err.Error())
					continue
				}

				if datatype == 3 {
					bufSlice = buf[resp.RESPONSE_HEADER_SIZE:]
					var partialRes []string

					err = json.Unmarshal(bufSlice, &partialRes)
					if err != nil {
						fmt.Printf("(err):%s\n", err.Error())
					}
					res = append(res, partialRes...)
				} else {
					fmt.Println("(err): invalid response type")
				}

			} else {
				res := buf[resp.RESPONSE_HEADER_SIZE:]
				l.Warn().Msgf("keys::forwardKeys %s", string(res))
			}
		}
	}
	return res
}
