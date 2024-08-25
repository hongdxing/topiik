/*
* author: duan hongxing
* date: 22 Jun 2024
* desc:
 */

package str

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

/*
* Set STRING key value, if key exists, value will be overrided, if TTL|TTLAT has value, TTL|TTLAT will also be overrided
* Parameter:
*	- pieces: command line that CMD(SET) stripped, the first piece is the KEY
* Return:
*	- OK if success
*	- SYNTAX_ERROR if syntax error
*
* Syntax: SET KEY VALUE [GET] [TTL seconds] [TTLAT unxix-time-seconds] [EX|NX]
*	- GET: Return old value of the key, or nil if key did not exist
 */
func Set(req datatype.Req) (result string, err error) {

	key := string(req.Keys[0])
	returnOld := false
	ttl := consts.INT64_MAX
	if len(req.Args) > 0 {
		pieces := strings.Split(req.Args, consts.SPACE)
		for i := 0; i < len(pieces); i++ {
			piece := strings.ToUpper(strings.TrimSpace(pieces[i]))
			if piece == "GET" {
				returnOld = true
			} else if strings.ToUpper(pieces[i]) == "TTL" {
				if len(pieces) <= i+1 {
					return "", errors.New(resp.RES_SYNTAX_ERROR + " near TTL")
				}
				ttl, err = strconv.ParseInt(pieces[i+1], 10, 64)
				if err != nil {
					return "", errors.New(resp.RES_SYNTAX_ERROR + " near TTL")
				}
				i++
				ttl += time.Now().UTC().Unix() // will overflow??? or should limit ttl user can type???
			} else if piece == "TTLAT" {
				if len(pieces) <= i+1 {
					return "", errors.New(resp.RES_SYNTAX_ERROR + " near TTLAT")
				}
				ttl, err = strconv.ParseInt(pieces[i+1], 10, 64)
				if err != nil {
					return "", errors.New(resp.RES_SYNTAX_ERROR + " near TTLAT")
				}
				i++
			} else if piece == "EX" {
				if _, ok := memo.MemMap[key]; ok {
					//
				} else {
					return "", errors.New(resp.RES_KEY_NOT_EXIST)
				}
			} else if piece == "NX" {
				if _, ok := memo.MemMap[key]; ok {
					return "", errors.New(resp.RES_KEY_EXIST_ALREADY)
				}
			} else {
				return "", errors.New(resp.RES_SYNTAX_ERROR)
			}
		}
	}

	if val, ok := memo.MemMap[key]; ok {
		/*
			if val.Typ != datatype.V_TYPE_STRING {
				return "", errors.New(resp.RES_DATA_TYPE_NOT_MATCH)
			}
		*/
		oldValue := resp.RES_OK
		if returnOld {
			oldValue = string(val.Str)
		}

		memo.MemMap[key].Str = []byte(req.Vals[0])
		memo.MemMap[key].Typ = datatype.V_TYPE_STRING
		memo.MemMap[key].Exp = ttl
		return oldValue, nil
	} else {
		memo.MemMap[key] = &datatype.TValue{
			Typ: datatype.V_TYPE_STRING,
			Str: []byte(req.Vals[0]),
			Exp: ttl}
		if returnOld {
			return "", nil
		}
		return resp.RES_OK, nil
	}
}
