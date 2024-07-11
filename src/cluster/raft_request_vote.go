/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"bufio"
	"fmt"
	"io"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"time"
	"topiik/internal/proto"
	"topiik/internal/util"
)

var voteMeResults []string // Return values from other Nodes, "R": Rejected or "A": Accepted
var wgRequestVote sync.WaitGroup

const requestVoteInterval = 100

/**
** Parameters
**	- addresses: workers' addresses
** Follower issues RequestVote RPCs to Workers to request for votes.
**/
func RequestVote() {

	// if this is the only Controller, then it alawys Leader
	if len(controllerMap) == 1 {
		nodeStatus.Role = RAFT_LEADER
		go AppendEntries()
		return
	}

	var quorum int
	//heartbeat := time.Duration(interval)
	var heartbeat time.Duration
	// Vote retry counter
	counter := 0
	for {
		quorum = 1 // I vote myself
		voteMeResults = voteMeResults[:0]

		heartbeat = time.Duration(rand.IntN(199)) + time.Duration(requestVoteInterval) //[0,99) + 200(interval), this must less than RaftHeartbeat(300)
		time.Sleep(heartbeat * time.Millisecond)

		if time.Now().UTC().UnixMilli() < nodeStatus.HeartbeatAt+int64(nodeStatus.Heartbeat) {
			if nodeStatus.Role != RAFT_FOLLOWER {
				nodeStatus.Role = RAFT_FOLLOWER
			}
			continue
		}
		// need workers to Vote
		/*if len(workerMap) == 0 {
			continue
		}*/
		// merge controller and woker address2
		var addr2List = []string{}
		for _, v := range controllerMap {
			addr2List = append(addr2List, v.Address2)
		}
		for _, v := range workerMap {
			addr2List = append(addr2List, v.Address2)
		}

		nodeStatus.Role = RAFT_CANDIDATOR
		nodeStatus.Term += 1
		for _, addr := range addr2List {
			wgRequestVote.Add(1)
			go voteMe(addr, int(nodeStatus.Term)) // use address2 for Voting
		}
		wgRequestVote.Wait()
		//fmt.Println(voteMeResults)
		for _, s := range voteMeResults {
			strs := strings.Split(s, ":")
			if len(strs) != 2 {
				break
			}
			if strings.Compare(strs[0], VOTE_REJECTED) == 0 {
				if strings.Compare(strs[1], "L") == 0 {
					//nodeStatus.Role = ROLE_FOLLOWER
					break
				}
			} else if strings.Compare(strs[0], "A") == 0 {
				quorum++
			}
		}

		canPromote := quorum >= ((len(addr2List))/2 + 1)
		if counter%10 == 0 || canPromote {
			fmt.Printf("Total nodes %v, quorum: %v\n", len(addr2List), quorum)
			// in case overflow
			if counter > 10000 {
				counter = 0
			}
		}
		if canPromote {
			// promote to Controller
			nodeStatus.Role = RAFT_LEADER
			// Leader start to AppendEntries
			go AppendEntries()
			// Print Selected Leader
			fmt.Printf(">>>selected as new Leader<<<\n")
			// Leader no RequestVote, quite RequestVote
			break
		} else {
			nodeStatus.Role = RAFT_FOLLOWER
		}
		counter++
	}
}

func voteMe(address string, term int) {
	defer wgRequestVote.Done()
	conn, err := util.PreapareSocketClient(address)
	if err != nil {
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

	reader := bufio.NewReader(conn)
	buf, err := proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("raft_request_vote::voteMe %s\n", err)
		}
	}
	voteMeResults = append(voteMeResults, string(buf[4:]))

	/*
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			return
		} else {
			//fmt.Println(string(buf))
			voteMeResults = append(voteMeResults, string(buf[:n]))
		}*/
}
