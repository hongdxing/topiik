package raft

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"topiik/internal/proto"
)

const (
	VOTE_REJECTED = "R"
	VOTE_ACCEPTED = "A"
)

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
func RequestVote(address *[]string, nodestatus *NodeStatus) {

	quota := 0
	heartbeat := time.Duration(1000 * 30)
	for {
		//time.Sleep(time.Duration(time.Duration(term).Milliseconds()))
		time.Sleep(heartbeat * time.Millisecond)
		nodestatus.Term += 1
		for _, addr := range *address {
			result := voteMe(addr, int(nodestatus.Term))
			if result == VOTE_ACCEPTED {
				quota++
			} else {

			}
		}

		if quota >= (len(*address) + 1) {
			// promote to Leader

			// Leader no RequestVote
			break
		}
	}
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

func voteMe(address string, term int) string {
	tcpServer, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())

	}
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println(err)
		return VOTE_REJECTED
	}

	line := "VOTE " + strconv.Itoa(term)

	// Enocde
	data, err := proto.Encode(line)
	if err != nil {
		fmt.Println(err)
	}

	// Send
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
		return VOTE_REJECTED
	}

	buf := make([]byte, 512)
	conn.Read(buf)
	fmt.Println(string(buf))
	return string(buf)

	//go response(conn)

}
