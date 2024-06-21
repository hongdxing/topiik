package raft

import (
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"
	"sync"
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

var voteMeResults []string
var wg sync.WaitGroup

/*
* Candidate issues RequestVote RPCs to other nodes to request for votes.
 */
func RequestVote(address *[]string, nodestatus *NodeStatus) {

	var quota int
	//var voteMeResults []string
	var heartbeat time.Duration
	for {
		quota = 0
		voteMeResults = voteMeResults[:0]
		heartbeat = time.Duration(rand.IntN(1000-500) + 500)
		time.Sleep(heartbeat * time.Millisecond)
		// Change role to Candidator
		nodestatus.Role = ROLE_CANDIDATOR
		nodestatus.Term += 1
		for _, addr := range *address {
			/*
				result := voteMe(addr, int(nodestatus.Term))
				fmt.Printf("voteMe result: %s \n", result)
				if result == VOTE_ACCEPTED {
					quota++
					fmt.Println(quota)
				} else {

				}*/
			wg.Add(1)
			go voteMe(addr, int(nodestatus.Term))
		}
		//fmt.Println(voteMeResults)
		for _, s := range voteMeResults {
			//fmt.Printf("----------%q---------\n", s)
			if strings.Compare(s, "A") == 0 {
				quota++
			}
		}

		fmt.Printf("Total nodes %v, quota: %v\n", len(*address)+1, quota)
		if quota >= ((len(*address)+1)/2 + 1) {
			// promote to Leader
			nodestatus.Role = ROLE_LEADER
			// Leader no RequestVote
			break
		} else {
			nodestatus.Role = ROLE_FOLLOWER
		}
	}
}

func voteMe(address string, term int) string {
	defer wg.Done()
	tcpServer, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())

	}
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println(err)
		return VOTE_REJECTED
	}
	defer conn.Close()

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
	n, err := conn.Read(buf)
	if err != nil {
		voteMeResults = append(voteMeResults, VOTE_REJECTED)
	} else {
		//fmt.Println(string(buf))
		voteMeResults = append(voteMeResults, string(buf[:n]))
	}
	return string(buf)

	//go response(conn)

}
