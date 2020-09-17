package raft

import (
	"log"
	"sort"
	"time"
)

func (rf *Raft) leaderRoutine() string {

	//log.Printf("%d leader,last index is %d,content is %v", rf.me, len(rf.status.Log)-1, rf.status.lastEntry())

	appendReplyChan := make(chan AppendEntryReply)

	rf.LeaderAppendHeatBeat()

	//init follower info
	followerInfos := make(map[int]FollowerInfo)
	for id := range rf.peers {
		if id == rf.me {
			continue
		}
		followerInfos[id] = FollowerInfo{rf.status.lastIndex(), 0}
	}
	rf.followersInfo = followerInfos
	rf.LeaderSyncLog(appendReplyChan)

	nextCommandIndex := 1
	lastCommandEntry, ok := rf.status.lastEntryContainsOperation()
	if ok {
		nextCommandIndex = lastCommandEntry.Command.Index + 1
	}
	timer := time.NewTimer(HeatBeatTimeout)
	for {
		select {
		// todo client Command
		case appendEntryArgs := <-rf.appendEntryRequest:
			// todo change to follower
			if rf.status.CurrentTerm == appendEntryArgs.CurrentTerm {
				log.Fatalf("%d two lead in same CurrentTerm", rf.me)
				break
			}

			res := rf.handleAppend(&appendEntryArgs)
			rf.appendEntryReply <- res
			if appendEntryArgs.CurrentTerm > rf.status.CurrentTerm {
				//log.Printf("leader %d receive append request ,update CurrentTerm from %d to %d", rf.me, rf.status.CurrentTerm,
				//	appendEntryArgs.CurrentTerm)
				rf.setCurrentTerm(appendEntryArgs.CurrentTerm)
				rf.persist()
				return follower
			}
		case voteArgs := <-rf.voteRequestChan:
			voteReply := rf.handleVote(voteArgs)
			rf.voteReplyChan <- voteReply

			if voteArgs.Term > rf.status.CurrentTerm {
				//log.Printf("leader %d receive a vote from %d, turn to follower",
				//	rf.me, voteArgs.Id)
				rf.setCurrentTerm(voteArgs.Term)
				rf.persist()
				return follower
			}

		case reply := <-appendReplyChan:
			id := reply.Id
			switch reply.Result {
			case appendEntryStaleTerm:
				//log.Printf("%d recieve more later CurrentTerm ,change from leader to follower", rf.me)
				rf.setCurrentTerm(reply.Term)
				return follower
			case appendEntryAccept:
				followerInfo := rf.followersInfo[id]
				followerInfo.NextIndex = reply.LastIndex + 1

				followerInfo.MatchIndex = followerInfo.NextIndex - 1
				//log.Printf("%d update %d next index to %d", rf.me, id, followerInfo.NextIndex)
				rf.followersInfo[id] = followerInfo
				rf.LeaderSyncCommittedIndex()
			case appendEntryNotMatch:
				followerInfo := rf.followersInfo[id]
				followerInfo.NextIndex -= appendConflictDecreaseNumber
				if followerInfo.NextIndex < 1 {
					followerInfo.NextIndex = 1
				}
				//log.Printf("%d decrease %d log index to %d", rf.me, id, followerInfo.NextIndex)
				rf.followersInfo[id] = followerInfo
			}

		case operation := <-rf.startRequestChan:
			command := Command{operation, nextCommandIndex}
			nextCommandIndex++
			lastEntry := rf.status.lastEntry()

			entry := Entry{rf.status.CurrentTerm, false, command, lastEntry.Index + 1}
			rf.status.Log = append(rf.status.Log, entry)
			rf.persist()
			rf.LeaderSyncLog(appendReplyChan)
			rf.startReplyChan <- StartReply{command.Index, int(rf.status.CurrentTerm)}

		case <-timer.C:
			if rf.killed() {
				log.Printf("%d is killed,exit", rf.me)
				return "exit"
			}
			//log.Printf("routine number is %d", runtime.NumGoroutine())
			rf.LeaderAppendHeatBeat()
			//log.Printf("%d send heart beat after time interval", rf.me)
			rf.LeaderSyncLog(appendReplyChan)
			timer = time.NewTimer(HeatBeatTimeout)

		}
	}
}

// append heat beat to Log
func (rf *Raft) LeaderAppendHeatBeat() {
	rf.status.Log = append(rf.status.Log, newEmptyEntry(rf.status.CurrentTerm, rf.status.lastIndex()+1))
	rf.persist()
}

func (rf *Raft) LeaderSyncCommittedIndex() {

	defer rf.persist()

	var followerMatchIndexList []int
	for _, value := range rf.followersInfo {
		followerMatchIndexList = append(followerMatchIndexList, int(value.MatchIndex))
	}
	sort.Ints(followerMatchIndexList)
	// add 1: include self
	committedIndex := Index(followerMatchIndexList[len(followerMatchIndexList)/2])
	if committedIndex > rf.committeeIndex {
		//log.Printf("%d update committee index from %d to %d", rf.me, rf.committeeIndex, committedIndex)
		for i := rf.committeeIndex + 1; i <= committedIndex; i++ {
			entry, _, _ := rf.status.getEntry(i)
			if entry.IsHeatBeat {
				continue
			}
			applyMsg := ApplyMsg{true, entry.Command.Operation,
				entry.Command.Index}
			//log.Printf("%d leader apply msg %v", rf.me, applyMsg)
			rf.ApplyMsgUnblockChan <- applyMsg
		}
		rf.committeeIndex = Index(committedIndex)
		rf.lastApply = rf.committeeIndex
	}
}

// sync leader Log to all peer
func (rf *Raft) LeaderSyncLog(appendReplyChan chan AppendEntryReply) {
	//log.Printf("%d begin sync log,current term %d", rf.me, rf.status.CurrentTerm)
	for id := range rf.peers {
		if id == rf.me {
			continue
		}
		followerInfo := rf.followersInfo[id]
		previousIndex := followerInfo.NextIndex - 1
		//log.Printf("follower %d info is %v", id, followerInfo)
		request := AppendEntryArgs{previousIndex, rf.status.getTerm(previousIndex),
			rf.status.Log[followerInfo.NextIndex:], rf.status.CurrentTerm,
			rf.committeeIndex}
		go func(rf *Raft, id int, args AppendEntryArgs) {
			reply := AppendEntryReply{}
			rf.sendAppendEntry(id, &args, &reply)
			timer := time.NewTimer(time.Second)
			for {
				select {
				case appendReplyChan <- reply:
					return
				case <-timer.C:
					return
				}
			}
		}(rf, id, request)
	}
}
