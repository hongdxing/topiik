//author: Duan Hongxing
//date: 11 Set, 2024

package list

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"topiik/internal/datatype"
	"topiik/memo"
	"topiik/resp"
)

// syntax: LSLICE Key start:end
// return:
//   - slice(subset) of list
//   - empty list if Key not exists
//
// arguments:
// start and end:
//   - "end" is exclusive
//   - both start and end are optional
//   - if start not specified, then start default to 0
//   - if end not specified, then end default lenght of the list
//   - if both start and end are not specified, then return whole list
//   - both start and end can be negative, -1 is the last element position
//
// note: lslice do not change the original list
func Slice(req datatype.Req) (rslt []string, err error) {
	if len(req.Args) == 0 {
		return nil, errors.New(resp.RES_SYNTAX_ERROR)
	}
	key := string(req.Keys[0])
	if val, ok := memo.MemMap[key]; ok {
		var (
			start = 0
			end   = val.Lst.Len()
		)
		se := strings.Split(strings.TrimSpace(req.Args), ":")
		fmt.Println(se)
		if len(se) == 2 {
			if len(strings.TrimSpace(se[0])) != 0 {
				start, err = strconv.Atoi(se[0])
				if err != nil {
					return nil, errors.New(resp.RES_SYNTAX_ERROR)
				}
			}
			if len(strings.TrimSpace(se[1])) != 0 {
				end, err = strconv.Atoi(se[1])
				if err != nil {
					return nil, errors.New(resp.RES_SYNTAX_ERROR)
				}
			}
		} else if req.Args[0] == ':' { //only specified end
			end, err = strconv.Atoi(se[1])
			if err != nil {
				return nil, errors.New(resp.RES_SYNTAX_ERROR)
			}
		} else if req.Args[len(req.Args)-1] == ':' { //only specified start
			start, err = strconv.Atoi(se[0])
			if err != nil {
				return nil, errors.New(resp.RES_SYNTAX_ERROR)
			}
		} else {
			return nil, errors.New(resp.RES_SYNTAX_ERROR)
		}
		if start < 0 {
			start = val.Lst.Len() + start
		}
		if end < 0 {
			end = val.Lst.Len() + end
		}

		// if start or end still negative after plus len, set to 0
		if start < 0 {
			start = 0
		}
		if end < 0 {
			end = 0
		}

		//fmt.Printf("start %v, end %v\n", start, end)

		looper := 0
		for ele := val.Lst.Front(); ele != nil && looper < end; ele = ele.Next() {
			if looper >= start {
				rslt = append(rslt, string(ele.Value.([]byte)))
			}
			looper++
		}
	}

	return rslt, err
}
