package raft

import (
	"log"
	"time"
)

func (rf *Raft) leaderRoutine() string {
	//args := rf.getAppendArgForHeatBeat()
	//reply := AppendEntryReply{}
	//for _, peer := range rf.peers {
	// todo go on here
	//todo handle reply async
	//go rf.LeaderSyncLog(peer,[]Command{newEmptyEntry()})
	//}

	// set to leader
	rf.isLeader = true
	defer func() { rf.isLeader = false }()

	rf.LeaderAppendHeatBeat()

	//init follower info
	followerInfos := make(map[int]FollowerInfo)
	for id := range rf.peers {
		followerInfos[id] = FollowerInfo{Index(len(rf.status.Log) - 1), 0}
	}
	rf.followersInfo = followerInfos

	rf.LeaderSyncLog()
	for {
		select {
		// todo client Command
		//case
		case appendEntryArgs := <-rf.appendEntryRequest:
			// todo change to follower
			if rf.status.CurrentTerm == appendEntryArgs.CurrentTerm {
				log.Fatalln("two lead in same CurrentTerm")
				break
			}

			res := rf.handleAppend(&appendEntryArgs)
			rf.appendEntryReply <- res
			if res.Result != appendEntryStaleTerm {
				log.Printf("leader %d accept a append request ,turn to follower", rf.me)
				return follower
			}
		case voteArgs := <-rf.voteRequestChan:
			voteReply := rf.handleVote(voteArgs)
			rf.voteReplyChan <- voteReply
			if voteReply.Result == voteReplySuccess {
				log.Printf("leader %d accept a vote from %d, turn to follower",
					rf.me, voteArgs.Id)
				return follower
			}
		case <-time.After(HeatBeatTimeout):
			rf.LeaderAppendHeatBeat()
			rf.LeaderSyncLog()
		}
	}
	return ""
}

// append heat beat to Log
func (rf *Raft) LeaderAppendHeatBeat() {
	rf.status.Log = append(rf.status.Log, newEmptyEntry(rf.status.CurrentTerm))

}

// sync leader Log to all peer
func (rf *Raft) LeaderSyncLog() {
	log.Printf("%d begin sync log", rf.me)
	for id := range rf.peers {
		if id == rf.me {
			continue
		}
		followerInfo := rf.followersInfo[id]
		previousIndex := followerInfo.nextIndex - 1
		request := AppendEntryArgs{previousIndex, rf.status.getTerm(previousIndex),
			rf.status.Log[followerInfo.nextIndex:], rf.status.CurrentTerm}
		go func(id int) {
			reply := AppendEntryReply{}
			rf.sendAppendEntry(id, &request, &reply)
		}(id)
	}
}
