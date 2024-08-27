/*
* author: duan hongxing
* date: 23 Jun 2024
* desc:
*
 */

package list

import (
	"fmt"
	"strconv"
	"strings"
	"topiik/internal/consts"
	"topiik/internal/datatype"
	"topiik/memo"
)

func Rang(req datatype.Req) (rslt []string, err error) {
	key := string(req.Keys[0])
	if value, ok := memo.MemMap[key]; ok {
		if value.Typ == datatype.V_TYPE_LIST && value.Lst.Len() > 0 {
			pieces := strings.Split(req.Args, consts.SPACE)
			if len(pieces) != 2 {
				return rslt, err
			}
			var start int
			var end int
			start, err = strconv.Atoi(pieces[0])
			if err != nil {
				return rslt, err
			}
			end, err = strconv.Atoi(pieces[1])
			if err != nil {
				return rslt, err
			}
			/*
			* [0...1...2...3...4]
			*  ^               ^
			*  |               |
			* start           end
			 */
			/* start gt right bound || end from 0 || end < 0 and end lt left bound */
			if start > value.Lst.Len() || end == 0 || (end < 0 && -end >= value.Lst.Len()) {
				return rslt, err
			}

			if start < 0 {
				if -start >= value.Lst.Len() {
					start = 0
				} else {
					start = value.Lst.Len() + start
				}
			}

			if end < 0 {
				/* end=-1; then end=5-1+1=5*/
				end = value.Lst.Len() + end + 1
			}
			fmt.Printf("start %v, end %v\n", start, end)

			looper := 0
			for ele := value.Lst.Front(); ele != nil && looper < end; ele = ele.Next() {
				if looper >= start {
					rslt = append(rslt, string(ele.Value.([]byte)))
				}
				looper++
			}
		}
	}

	return rslt, err
}
