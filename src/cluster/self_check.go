/***
** author: duan hongxing
** date: 30 Jun 2024
** desc:
**	return keys
**
**/

package cluster

/***
** When node up
	1) The role always FOLLOW
	2) If /data/cluster file NOT exists
		- If no cluster file, then Single node mode, untill join a cluster
		-
	3) If ./data/cluster file exists, and 3 nodes mode
		- Send request to either of the other nodes to rejoin the cluster
		- If accepted by other nodes, change role to MASTER
		- If rejected by other nodes, then return error
	4) If ./data/cluster file exists, and 6 nodes mode

**
**
**
**
**/
func SelfCheck() {

}
