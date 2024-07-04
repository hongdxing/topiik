/***
** author: duan hongxing
** data: 4 Jul 2024
** desc:
**
**/

package executer

import "topiik/ccss"

func clusterJoinACK(pieces []string)(result string, err error) {
	return ccss.JoinACK(pieces[0], pieces[1])
}
