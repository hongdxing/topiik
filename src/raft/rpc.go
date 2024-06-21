package raft

/*
* author: duan hongxing
* date: 21 Jun 2024
* desc: Raft RPC implements
* ref:
*	- https://www.sofastack.tech/en/projects/sofa-jraft/consistency-raft-jraft/
*
 */

/*
* Candidate issues RequestVote RPCs to other nodes to request for votes.
 */
func RequestVote() {

}

/***
* leader issues AppendEntries RPCs to replicate log entries to followers,
* or send heartbeats (AppendEntries RPCs that carry no log entries)
 */
func AppendEntries() {

}

/***
* InstallSnapshot
 */
func InstallSnapshot() {

}
