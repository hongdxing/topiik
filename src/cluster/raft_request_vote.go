//author: duan hongxing
//data: 3 Jul 2024

package cluster

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"time"
	"topiik/internal/consts"
	"topiik/internal/proto"
	"topiik/internal/util"
	"topiik/node"
	"topiik/persistence"
	"topiik/resp"
)

var voteMeResults []string // Return values from other Nodes, "R": Rejected or "A": Accepted
var wgRequestVote sync.WaitGroup

const requestVoteInterval = 100

// Request vote self to controller leader
// Parameters:
// Follower(s) issue RequestVote RPCs to Controller(s) and Worker(s) to request for votes.
func RequestVote() {
	if nodeStatus.Role == RAFT_LEADER {
		return
	}
	wrkGrp := getWorkGroup(node.GetNodeInfo().Id)
	if wrkGrp.Id == "" {
		l.Warn().Msgf("RequestVote Invalid worker group")
		return
	}

	if len(wrkGrp.Nodes) == 1 { // if this is the only worker in group, then it alawys Leader
		nodeStatus.Role = RAFT_LEADER
		go AppendEntries(wrkGrp)
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

		heartbeat = time.Duration(rand.IntN(199)) + time.Duration(requestVoteInterval) //[0,199) + 100(interval), this must less than RaftHeartbeat(300)
		time.Sleep(heartbeat * time.Millisecond)

		if time.Now().UTC().UnixMilli() < nodeStatus.HeartbeatAt+int64(nodeStatus.Heartbeat) {
			if nodeStatus.Role != RAFT_FOLLOWER {
				nodeStatus.Role = RAFT_FOLLOWER
			}
			continue
		}
		// merge worker and persistor
		var addr2List = []string{}
		for _, v := range wrkGrp.Nodes {
			addr2List = append(addr2List, v.Addr2)
		}

		// use persistor as voter
		// todo: what if Persistor[0] down???
		pstAddr := node.GetConfig().Persistors[0]
		host, _, port2, _ := util.SplitAddress2(pstAddr)
		addr2List = append(addr2List, host+":"+port2)

		nodeStatus.Role = RAFT_CANDIDATOR
		nodeStatus.Term += 1
		for _, addr := range addr2List {
			wgRequestVote.Add(1)
			go voteMe(addr) // use address2 for Voting
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
					nodeStatus.Role = RAFT_FOLLOWER
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
			// promote to Leader
			nodeStatus.Role = RAFT_LEADER

			// Leader start to AppendEntries
			go AppendEntries(wrkGrp)
			go persistence.Dequeue()
			/* make sure channel are ready */
			//time.Sleep(500 * time.Millisecond)

			// when new Leader selected, try to sync cluster metadata
			notifyWorkerGroupChanged()

			/* Print Selected Leader */
			l.Info().Msgf("[TOPIIK] ~!~ selected as new leader")

			/* Leader no RequestVote, quite RequestVote */
			break
		} else {
			nodeStatus.Role = RAFT_FOLLOWER
		}
		counter++
	}
}

func voteMe(address string) {
	defer wgRequestVote.Done()
	conn, err := util.PreapareSocketClient(address)
	if err != nil {
		return
	}
	defer conn.Close()

	var buf []byte
	var bbuf = new(bytes.Buffer) // int to byte buf
	_ = binary.Write(bbuf, binary.LittleEndian, consts.RPC_VOTE)
	buf = append(buf, bbuf.Bytes()...)
	buf = append(buf, []byte(node.GetNodeInfo().Id)...) // Include Cluster Id in request
	buf = append(buf, []byte(strconv.Itoa(term))...)

	// Enocde
	buf, err = proto.EncodeB(buf)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}

	// Send
	_, err = conn.Write(buf)
	if err != nil {
		l.Err(err).Msg(err.Error())
		return
	}

	reader := bufio.NewReader(conn)
	buf, err = proto.Decode(reader)
	if err != nil {
		if err == io.EOF {
			l.Err(err).Msgf("raft_request_vote::voteMe %s", err)
		}
		return
	}
	if len(buf) < resp.RESPONSE_HEADER_SIZE {
		l.Err(err).Msgf("raft_request_vote::voteMe invalid len(buf):%v", len(buf))
		return
	}

	// If result is REJECTED, means current Node has been removed from cluster
	if string(buf[resp.RESPONSE_HEADER_SIZE:]) == resp.RES_REJECTED {
		l.Panic().Msg("Node rejected by other nodes, possible has been removed from cluster")
	}
	voteMeResults = append(voteMeResults, string(buf[resp.RESPONSE_HEADER_SIZE:]))
}
