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
*	SET: {"Keys": [k1], "Vals": [v1], "Args":"" }
* 	GETM: {"Keys": [k1 k2], "Vals": [], "Args":"" }
*	LPUSH: {"Keys": [list], "Vals": [111 222 333], "Args":"" }
 */

package datatype

/*
type Req struct {
	Keys []string
	Vals []string
	Args string
}
*/

/* Array of bytes */
type Abytes [][]byte

type Req struct {
	Keys Abytes
	Vals Abytes
	Args string
}
