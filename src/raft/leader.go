package raft

import (
	"log"
	"sort"
	"time"
)

func (rf *Raft) leaderRoutine() string {
	//Args := rf.getAppendArgForHeatBeat()
	//Reply := AppendEntryReply{}
	//for _, peer := range rf.peers {
	// todo go on here
	//todo handle Reply async
	//go rf.LeaderSyncLog(peer,[]Command{newEmptyEntry()})
	//}

	appendReplyChan := make(chan AppendEntryInfo)

	rf.LeaderAppendHeatBeat()

	//init follower info
	followerInfos := make(map[int]FollowerInfo)
	for id := range rf.peers {
		followerInfos[id] = FollowerInfo{Index(len(rf.status.Log) - 1), 0}
	}
	rf.followersInfo = followerInfos
	rf.LeaderSyncLog(appendReplyChan)

	// find latest committee command
	rf.nextCommandIndex = 1
	for i := len(rf.status.Log) - 1; i >= 0; i-- {
		//for _, entry := range rf.status.Log {
		entry := rf.status.Log[i]
		if !entry.IsHeatBeat {
			rf.nextCommandIndex = entry.Index + 1
			break
		}
	}

	for {
		timer := time.NewTimer(rf.randomElectionTimeout())
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

			if voteArgs.Term > rf.status.CurrentTerm {
				log.Printf("leader %d receive a vote from %d, turn to follower",
					rf.me, voteArgs.Id)
				rf.status.CurrentTerm = voteArgs.Term
				return follower
			}

		case appendInfo := <-appendReplyChan:
			id, request, reply := appendInfo.Id, appendInfo.Args, appendInfo.Reply
			switch reply.Result {
			case appendEntryStaleTerm:
				log.Printf("%d recieve more later CurrentTerm ,change from leader to follower", rf.me)
				rf.status.CurrentTerm = reply.Term
				return follower
			case appendEntryAccept:
				followerInfo := rf.followersInfo[id]
				followerInfo.NextIndex = Index(int(request.PreviousEntryIndex) + len(request.Entries) + 1)
				followerInfo.MatchIndex = followerInfo.NextIndex - 1
				rf.followersInfo[id] = followerInfo
				rf.LeaderSyncCommittedIndex()
			case appendEntryNotMatch:
				followerInfo := rf.followersInfo[id]
				followerInfo.NextIndex -= appendConflictDecreaseNumber
				if followerInfo.NextIndex < 1 {
					followerInfo.NextIndex = 1
				}
				rf.followersInfo[id] = followerInfo
			}

		case command := <-rf.startRequestChan:
			log.Printf("debug comand log index is %d", len(rf.status.Log)-1)
			entry := Entry{rf.status.CurrentTerm, false, command, rf.nextCommandIndex}
			rf.nextCommandIndex++
			rf.status.Log = append(rf.status.Log, entry)
			//rf.status.lastCommitteeIndex++
			rf.LeaderSyncLog(appendReplyChan)
			rf.startReplyChan <- StartReply{int(entry.Index), int(rf.status.CurrentTerm)}

		case <-timer.C:
			rf.LeaderAppendHeatBeat()
			log.Printf("%d send heart beat after time interval", rf.me)
			rf.LeaderSyncLog(appendReplyChan)
		}
	}
	return ""
}

// append heat beat to Log
func (rf *Raft) LeaderAppendHeatBeat() {
	rf.status.Log = append(rf.status.Log, newEmptyEntry(rf.status.CurrentTerm))

}

func (rf *Raft) LeaderSyncCommittedIndex() {
	var followerMatchIndexList []int
	for _, value := range rf.followersInfo {
		followerMatchIndexList = append(followerMatchIndexList, int(value.MatchIndex))
	}
	sort.Ints(followerMatchIndexList)
	// add 1: include self
	committedIndex := Index(followerMatchIndexList[len(followerMatchIndexList)/2+1])
	if committedIndex > rf.committeeIndex {
		log.Printf("%d update committee index from %d to %d", rf.me, rf.committeeIndex, committedIndex)
		for _, entry := range rf.status.Log[rf.committeeIndex+1 : committedIndex+1] {
			if entry.IsHeatBeat {
				continue
			}
			applyMsg := ApplyMsg{true, entry.Command.Content, int(entry.Index)}
			log.Printf("%d leader apply msg %v", rf.me, applyMsg)
			rf.applyMsgChan <- applyMsg
			rf.nextApplyIndex++
		}
		rf.committeeIndex = Index(committedIndex)
		//todo  apply to state machine
		rf.lastApply = rf.committeeIndex
	}
}

// sync leader Log to all peer
func (rf *Raft) LeaderSyncLog(appendReplyChan chan AppendEntryInfo) {
	log.Printf("%d begin sync log", rf.me)
	for id := range rf.peers {
		if id == rf.me {
			continue
		}
		followerInfo := rf.followersInfo[id]
		previousIndex := followerInfo.NextIndex - 1
		request := AppendEntryArgs{previousIndex, rf.status.getTerm(previousIndex),
			rf.status.Log[followerInfo.NextIndex:], rf.status.CurrentTerm,
			rf.committeeIndex}
		go func(id int) {
			reply := AppendEntryReply{}
			rf.sendAppendEntry(id, &request, &reply)
			appendReplyChan <- AppendEntryInfo{id, request, reply}
		}(id)
	}
}
