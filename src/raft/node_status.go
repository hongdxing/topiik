/***
* author: duan hongxing
* date: 21 Jun 2024
* desc: Store current node status
*
 */
package raft

type NodeStatus struct {
	Role        uint8  // Follower, Candidator, Leader
	Term        uint   // Raft term
	Heartbeat   uint16 // Raft heartbeat timeout
	HeartbeatAt int64  // The UTC milli seconds when heartbeat received from Leader
}
