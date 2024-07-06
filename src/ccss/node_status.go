/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package ccss

type NodeStatus struct {
	Role        uint8  // Captial, Chief, Sailor
	Term        uint   // Raft term
	Heartbeat   uint16 // Raft heartbeat timeout
	HeartbeatAt int64  // The UTC milli seconds when heartbeat received from Leader
}
