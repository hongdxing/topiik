/*
** author: duan hongxing
** date: 22 Jul 2024
** desc:
**	- struct of Request
**	- client build Req, marshal to json, and encode to bytes
**	- server receive the bytes, unmarshal to Req
**	- use struct will add some more overload on network compare to use raw string command
**	- the reason why use struct is, it will save time on parse raw string command
** sample:
**	{"VER": 1, "CMD": "SET", "KEYS": [k1], "VALS": [v1], "ARGS":"" }
** 	{"VER": 1, "CMD": "GETM","KEYS": [k1 k2], "VALS": [], "ARGS":"" }
**	{"VER": 1, "CMD": "LPUSH", "KEYS": [list], "VALS": [111 222 333], "ARGS":"" }
 */

package datatype

type Req struct {
	//VER  uint8
	//CMD  string
	KEYS []string
	VALS []string
	ARGS string
}
