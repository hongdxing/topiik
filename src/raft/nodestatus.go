package raft

type NodeStatus struct {
	Role      uint8
	Term      uint
	Heartbeat uint8
}
