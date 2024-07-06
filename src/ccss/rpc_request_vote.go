/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package ccss

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"topiik/internal/proto"
)

var voteMeResults []string // Return values from other Nodes, "R": Rejected or "A": Accepted
var wgRequestVote sync.WaitGroup

const requestVoteInterval = 200

/**
** Parameters
**	- addresses: salors' addresses
** Chief issues RequestVote RPCs to Salors to request for votes.
**/
func RequestVote() {

	var quorum int
	//heartbeat := time.Duration(interval)
	var heartbeat time.Duration
	// Vote retry counter
	counter := 0
	for {
		quorum = 1 // Initial value 1, means vote current node(I vote myself)
		voteMeResults = voteMeResults[:0]

		heartbeat = time.Duration(99) + time.Duration(requestVoteInterval) //[0,99) + 200(interval), this must less than RaftHeartbeat(300)
		time.Sleep(heartbeat * time.Millisecond)

		if len(salorMap) == 0 { // if no Salor, then no RequestVote
			continue
		}

		if time.Now().UTC().UnixMilli() < nodeStatus.HeartbeatAt+int64(nodeStatus.Heartbeat) {
			if nodeStatus.Role != CCSS_ROLE_CO {
				nodeStatus.Role = CCSS_ROLE_CO
			}
			continue
		}

		nodeStatus.Term += 1
		for _, salor := range salorMap {
			wgRequestVote.Add(1)
			go voteMe(salor.Address, int(nodeStatus.Term))
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
					//nodeStatus.Role = ROLE_FOLLOWER
					break
				}
			} else if strings.Compare(strs[0], "A") == 0 {
				quorum++
			}
		}

		canPromote := quorum >= ((len(salorMap)+1)/2 + 1)
		if counter%10 == 0 || canPromote {
			fmt.Printf("Total nodes %v, quota: %v\n", len(salorMap)+1, quorum)
			// in case overflow
			if counter > 10000 {
				counter = 0
			}
		}
		if canPromote {
			// promote to Capital
			nodeStatus.Role = CCSS_ROLE_CA
			// Leader start to AppendEntries
			//go AppendEntries(salorMap)
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
