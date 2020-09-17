package raft

import (
	"math"
	"time"
)

// return true if append success
func (rf *Raft) handleAppend(appendArg *AppendEntryArgs) AppendEntryReply {

	defer rf.persist()

	if rf.status.CurrentTerm > appendArg.CurrentTerm {
		//log.Printf("%d append entry's CurrentTerm mistach, current CurrentTerm is %d, "+
		//	"CurrentTerm in requeset is %d ,reject it", rf.me, rf.status.CurrentTerm, appendArg.CurrentTerm)
		return AppendEntryReply{rf.me, appendEntryStaleTerm, rf.status.CurrentTerm, -1}
	}

	//log.Printf("%d receive a append,previous Index is %d,term is %d, service last Index is %d,last entry term is %d",
	//	rf.me, appendArg.PreviousEntryIndex, appendArg.PreviousEntryTerm, len(rf.status.Log)-1, rf.status.lastEntry().Term)

	match := rf.status.logContain(appendArg.PreviousEntryIndex, appendArg.PreviousEntryTerm)

	if !match {
		//log.Printf("append entry prev entry match fail")
		return AppendEntryReply{rf.me, appendEntryNotMatch, rf.status.CurrentTerm, -1}
	}

	for i, newEntry := range appendArg.Entries {
		if entry, position, ok := rf.status.getEntry(newEntry.Index); ok {
			if entry.Term != newEntry.Term {
				rf.status.deleteAfter(position)
				rf.status.append(appendArg.Entries[i:])
			}
		} else {
			rf.status.append(appendArg.Entries[i:])
			break
			//rf.status.Log = append(rf.status.Log, newEntry)
		}
	}

	if appendArg.LeaderCommittee > rf.committeeIndex {
		for i := rf.committeeIndex + 1; i <= appendArg.LeaderCommittee; i++ {
			//for _, entry := range rf.status.Log[rf.committeeIndex+1 : appendArg.LeaderCommittee+1] {
			entry, _, ok := rf.status.getEntry(Index(i))

			if !ok {
				break
			}

			if entry.IsHeatBeat {
				continue
			}
			applyMsg := ApplyMsg{true, entry.Command.Operation,
				entry.Command.Index}
			rf.ApplyMsgUnblockChan <- applyMsg
		}

		min := Index(math.Min(float64(appendArg.LeaderCommittee), float64(rf.status.lastIndex())))
		//log.Printf("%d update committee Index from %d to %d", rf.me, rf.committeeIndex, min)
		rf.committeeIndex = min
		rf.lastApply = min
	}
	//log.Printf("append finish, log size is %d", len(rf.status.Log))
	return AppendEntryReply{rf.me, appendEntryAccept,
		rf.status.CurrentTerm, rf.status.lastIndex()}

}

func (rf *Raft) handleVote(voteArgs RequestVoteArgs) (reply RequestVoteReply) {

	reply.CurrentTerm = rf.status.CurrentTerm
	if voteArgs.Term <= rf.status.CurrentTerm {
		//log.Printf("%d receive vote request from an old CurrentTerm, current CurrentTerm is %d, Reply CurrentTerm is %d, refuse it",
		//	rf.me, rf.status.CurrentTerm, voteArgs.Term)
		reply.Result = voteReplyStaleTerm
		return
	}

	// compare latest entry
	lastSlotIndex := rf.status.lastIndex()
	lastSlotTerm := rf.status.Log[lastSlotIndex].Term

	//log.Printf("%d receive vote, service Index is %d,term is %d,vote Id is %d ,Index is %d,term is %d",
	//	rf.me, lastSlotIndex, lastSlotTerm, voteArgs.Id, voteArgs.LastSlotIndex, voteArgs.LastSlotTerm)

	if lastSlotTerm > voteArgs.LastSlotTerm ||
		(lastSlotTerm == voteArgs.LastSlotTerm && lastSlotIndex > voteArgs.LastSlotIndex) {
		//log.Printf("%d last slot CurrentTerm  is %d, last Index is %d "+
		//	"later than vote's CurrentTerm %d, Index %d reject vote",
		//	rf.me, lastSlotTerm, lastSlotIndex, voteArgs.LastSlotTerm, voteArgs.LastSlotIndex)
		reply.Result = voteReplyLatestLogEntryIsNotUpdateToMe
		return
	}

	if _, ok := rf.status.VoteFor[voteArgs.Term]; ok {
		//log.Printf("%d CurrentTerm %d is already vote", rf.me, voteArgs.Term)
		reply.Result = voteReplyApplyAlreadyVote
		return
	}

	//log.Printf("%d vote for service %d", rf.me, voteArgs.Id)

	rf.status.VoteFor[voteArgs.Term] = voteArgs.Id
	rf.persist()
	reply.Result = voteReplySuccess
	return
}

func (rf *Raft) applyMsgRoutine() {
	for {
		timer := time.NewTimer(time.Second * 5)
		select {
		case applyMsg := <-rf.ApplyMsgUnblockChan:
			rf.applyMsgChan <- applyMsg
		case <-timer.C:
			timer = time.NewTimer(time.Second * 5)
			if rf.killed() {
				return
			}
		}
	}
}

func (rf *Raft) isLeader() bool {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.leaderStatus

}

func (rf *Raft) changeLeader(leaderStatus bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.leaderStatus = leaderStatus
}

func (rf *Raft) getCurrentTerm() Term {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.status.CurrentTerm
}

func (rf *Raft) incCurrentTerm() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.status.CurrentTerm++
}
func (rf *Raft) setCurrentTerm(term Term) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.status.CurrentTerm = term
}
