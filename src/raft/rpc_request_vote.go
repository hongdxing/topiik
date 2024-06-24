package raft

import (
	"fmt"
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

var voteMeResults []string // Return values from other Nodes, "R": Rejected or "A": Accepted
var wgRequestVote sync.WaitGroup

/*
* Candidate issues RequestVote RPCs to other nodes to request for votes.
 */
func RequestVote(addresses *[]string, interval uint16, nodestatus *NodeStatus) {

	var quorum int
	//heartbeat := time.Duration(interval)
	var heartbeat time.Duration
	// Vote retry counter
	counter := 0
	for {
		quorum = 1 // Initial value 1, means vote current node(I vote myself)
		voteMeResults = voteMeResults[:0]

		heartbeat = time.Duration(99) + time.Duration(interval) //[0,99) + 200(interval), this must less than RaftHeartbeat(300)
		time.Sleep(heartbeat * time.Millisecond)
		if time.Now().UTC().UnixMilli() < nodestatus.HeartbeatAt+int64(nodestatus.Heartbeat) {
			if nodestatus.Role != ROLE_FOLLOWER {
				nodestatus.Role = ROLE_FOLLOWER
			}
			continue
		}

		// Change role to Candidator
		nodestatus.Role = ROLE_CANDIDATOR
		nodestatus.Term += 1
		for _, addr := range *addresses {
			wgRequestVote.Add(1)
			go voteMe(addr, int(nodestatus.Term))
		}
		//fmt.Println(voteMeResults)
		for _, s := range voteMeResults {
			//fmt.Printf("----------%q---------\n", s)
			strs := strings.Split(s, ":")
			if len(strs) != 2 {
				break
			}
			if strings.Compare(strs[0], VOTE_REJECTED) == 0 {
				if strings.Compare(strs[1], "L") == 0 {
					//nodestatus.Role = ROLE_FOLLOWER
					break
				}
			} else if strings.Compare(strs[0], "A") == 0 {
				quorum++
			}
		}

		canPromote := quorum >= ((len(*addresses)+1)/2 + 1)
		if counter%10 == 0 || canPromote {
			fmt.Printf("Total nodes %v, quota: %v\n", len(*addresses)+1, quorum)
			// in case overflow
			if counter > 10000 {
				counter = 0
			}
		}
		if canPromote {
			// promote to Leader
			nodestatus.Role = ROLE_LEADER
			// Leader start to AppendEntries
			go AppendEntries(*addresses)
			// Print Selected Leader
			fmt.Printf(">>>selected as new Leader<<<\n")
			// Leader no RequestVote, quite RequestVote
			break
		} else {
			//nodestatus.Role = ROLE_FOLLOWER
		}
		counter++
	}
}

func voteMe(address string, term int) {
	defer wgRequestVote.Done()
	tcpServer, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		//fmt.Println(err)
		return
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
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return
	} else {
		//fmt.Println(string(buf))
		voteMeResults = append(voteMeResults, string(buf[:n]))
	}
}
