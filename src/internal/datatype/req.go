/*
* author: duan hongxing
* date: 22 Jul 2024
* desc:
*	- struct of Request
*	- client build Req, marshal to json, and encode to bytes
*	- server receive the bytes, unmarshal to Req
*	- use struct will add some more overload on network compare to use raw string command
*	- the reason why use struct is, it will save time on parse raw string command
* sample:
*	SET: {"KEYS": [k1], "VALS": [v1], "ARGS":"" }
* 	GETM: {"KEYS": [k1 k2], "VALS": [], "ARGS":"" }
*	LPUSH: {"KEYS": [list], "VALS": [111 222 333], "ARGS":"" }
 */

package datatype

type Req struct {
	KEYS []string
	VALS []string
	ARGS string
}
