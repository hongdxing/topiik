/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package cluster

import (
	"fmt"
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

const requestVoteInterval = 200

/**
** Parameters
**	- addresses: workers' addresses
** Chief issues RequestVote RPCs to Workers to request for votes.
**/
func RequestVote() {

	var quorum int
	//heartbeat := time.Duration(interval)
	var heartbeat time.Duration
	// Vote retry counter
	counter := 0
	for {
		quorum = 0 //
		voteMeResults = voteMeResults[:0]

		heartbeat = time.Duration(rand.IntN(99)) + time.Duration(requestVoteInterval) //[0,99) + 200(interval), this must less than RaftHeartbeat(300)
		time.Sleep(heartbeat * time.Millisecond)

		if time.Now().UTC().UnixMilli() < nodeStatus.HeartbeatAt+int64(nodeStatus.Heartbeat) {
			if nodeStatus.Role != RAFT_FOLLOWER {
				nodeStatus.Role = RAFT_FOLLOWER
			}
			continue
		}
		// need workers to Vote
		if len(workerMap) == 0 {
			continue
		}

		nodeStatus.Role = RAFT_CANDIDATOR
		nodeStatus.Term += 1
		for _, worker := range workerMap {
			wgRequestVote.Add(1)
			go voteMe(worker.Address2, int(nodeStatus.Term)) // use address2 for Voting
		}
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

		canPromote := quorum >= ((len(workerMap))/2 + 1)
		if counter%10 == 0 || canPromote {
			fmt.Printf("Total nodes %v, quorum: %v\n", len(workerMap)+1, quorum)
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
			//nodeStatus.Role = ROLE_FOLLOWER
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

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return
	} else {
		//fmt.Println(string(buf))
		voteMeResults = append(voteMeResults, string(buf[:n]))
	}
}
