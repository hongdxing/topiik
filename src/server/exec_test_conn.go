// author: Duan HongXing
// date: 7 Sep, 2024

package server

import "topiik/node"

func testConn() string {
	return node.GetNodeInfo().Id
}
