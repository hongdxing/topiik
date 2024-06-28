/***
* author: duan hongxing
* date: 22 Jun 2024
* desc:
***/

package executor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/shared"
)

/***
** Set STRING key value, if key exists, value will be overrided, if TTL|TTLAT has value, TTL|TTLAT will also be overrided
** Parameter:
**	- pieces: command line that CMD(SET) stripped, the first piece is the KEY
** Return:
**	- OK if success
**	- SYNTAX_ERROR if syntax error
**
** Syntax: SET KEY VALUE [GET] [TTL seconds] | [TTLAT unxix-time-seconds] | [EX|NX]
**	- GET: Return old value of the key, or nil if key did not exist
**/
func set(pieces []string) (result any, err error) {
	if len(pieces) < 2 {
		return "", errors.New(RES_SYNTAX_ERROR)
	}

	key := strings.TrimSpace(pieces[0])
	returnOld := false
	ttl := consts.UINT32_MAX
	if len(pieces) > 2 {
		for i := 2; i < len(pieces); i++ {
			if i == 0 || i == 1 {
				continue // skip KEY VALUE
			}
			piece := strings.ToUpper(strings.TrimSpace(pieces[i]))
			fmt.Printf("---%s---\n", piece)
			if piece == "GET" {
				returnOld = true
			} else if strings.ToUpper(pieces[i]) == "TTL" {
				if len(pieces) <= i+1 {
					return nil, errors.New(RES_SYNTAX_ERROR + " near TTL")
				}
				ttl, err = strconv.Atoi(pieces[i+1])
				if err != nil {
					return nil, errors.New(RES_SYNTAX_ERROR + " near TTL")
				}
				i++
				ttl += time.Now().UTC().Second() // will overflow??? or should limit ttl user can type???
			} else if piece == "TTLAT" {
				if len(pieces) <= i+1 {
					return nil, errors.New(RES_SYNTAX_ERROR + " near TTLAT")
				}
				ttl, err = strconv.Atoi(pieces[i+1])
				if err != nil {
					return nil, errors.New(RES_SYNTAX_ERROR + " near TTLAT")
				}
				i++
			} else if piece == "EX" {
				fmt.Println("EX")
				if _, ok := shared.MemMap[key]; ok {
					//
				} else {
					return nil, errors.New(RES_KEY_NOT_EXIST)
				}
			} else if piece == "NX" {
				fmt.Println("NX")
				if _, ok := shared.MemMap[key]; ok {
					return nil, errors.New(RES_KEY_EXIST_ALREADY)
				}
			}
		}
	}

	if val, ok := shared.MemMap[key]; ok {
		if val.Typ != datatype.V_TYPE_STRING {
			return "", errors.New(RES_DATA_TYPE_NOT_MATCH)
		}
		oldValue := RES_OK
		if returnOld {
			oldValue = string(val.Str)
		}

		shared.MemMap[key].Str = []byte(pieces[1])
		shared.MemMap[key].Exp = uint32(ttl)
		return oldValue, nil
	} else {
		shared.MemMap[key] = &datatype.TValue{
			Typ: datatype.V_TYPE_STRING,
			Str: []byte(pieces[1]),
			Exp: uint32(ttl)}
		if returnOld {
			return nil, nil
		}
		return RES_OK, nil
	}
}
