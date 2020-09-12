package raft

import (
	"log"
	"math"
)

//return Log index  if match
func (s *Status) logContain(index Index, term Term) bool {
	return len(s.Log) > int(index) && s.Log[index].Term == term
}

// return true if append success
func (rf *Raft) handleAppend(appendArg *AppendEntryArgs) AppendEntryReply {

	if rf.status.CurrentTerm > appendArg.CurrentTerm {
		log.Printf("%d append entry's CurrentTerm mistach, current CurrentTerm is %d, "+
			"CurrentTerm in requeset is %d ,reject it", rf.me, rf.status.CurrentTerm, appendArg.CurrentTerm)
		return AppendEntryReply{rf.me, appendEntryStaleTerm, rf.status.CurrentTerm, -1}
	}

	log.Printf("%d receive a append,previous index is %d,term is %d, service last index is %d,last entry term is %d",
		rf.me, appendArg.PreviousEntryIndex, appendArg.PreviousEntryTerm, len(rf.status.Log)-1, rf.status.lastEntry().Term)

	match := rf.status.logContain(appendArg.PreviousEntryIndex, appendArg.PreviousEntryTerm)

	if !match {
		log.Printf("append entry prev entry match fail")
		return AppendEntryReply{rf.me, appendEntryNotMatch, rf.status.CurrentTerm, -1}
	}

	// delete log conflict
	//if rf.status.Log[int(appendArg.PreviousEntryIndex)].Term != appendArg.PreviousEntryTerm {
	//	log.Printf("%d delete all Log after the match,from %d to %d",
	//		rf.me, appendArg.PreviousEntryIndex+1, len(rf.status.Log)-1)
	//	rf.status.Log = rf.status.Log[:appendArg.PreviousEntryIndex+1]
	//}

	for i := 0; i < len(appendArg.Entries); i++ {
		newEntry := appendArg.Entries[i]
		position := int(appendArg.PreviousEntryIndex) + 1 + i
		if position < len(rf.status.Log) {
			entry := rf.status.Log[position]
			if entry.Term != newEntry.Term {
				rf.status.Log[position] = appendArg.Entries[i]
				rf.status.Log = rf.status.Log[:position+1]
			}
		} else {
			rf.status.Log = append(rf.status.Log, newEntry)
		}
	}

	if appendArg.LeaderCommittee > rf.committeeIndex {
		//todo  apply to state machine
		for _, entry := range rf.status.Log[rf.committeeIndex+1 : appendArg.LeaderCommittee+1] {
			if entry.IsHeatBeat {
				continue
			}
			applyMsg := ApplyMsg{true, entry.Command.Content,
				int(entry.Index)}
			rf.applyMsgChan <- applyMsg
		}

		min := Index(math.Min(float64(appendArg.LeaderCommittee), float64(len(rf.status.Log)-1)))
		log.Printf("%d update committee index from %d to %d", rf.me, rf.committeeIndex, min)
		rf.committeeIndex = min
		rf.lastApply = min
	}
	rf.persist()
	log.Printf("append finish, log size is %d", len(rf.status.Log))
	return AppendEntryReply{rf.me, appendEntryAccept,
		rf.status.CurrentTerm, Index(len(rf.status.Log) - 1)}

}

func (rf *Raft) handleVote(voteArgs RequestVoteArgs) (reply RequestVoteReply) {

	reply.CurrentTerm = rf.status.CurrentTerm
	if voteArgs.Term <= rf.status.CurrentTerm {
		log.Printf("%d receive vote request from an old CurrentTerm, current CurrentTerm is %d, Reply CurrentTerm is %d, refuse it",
			rf.me, rf.status.CurrentTerm, voteArgs.Term)
		reply.Result = voteReplyStaleTerm
		return
	}

	// compare latest entry
	lastSlotIndex := len(rf.status.Log) - 1
	lastSlotTerm := rf.status.Log[lastSlotIndex].Term

	log.Printf("%d receive vote, service index is %d,term is %d,vote Id is %d ,index is %d,term is %d",
		rf.me, lastSlotIndex, lastSlotTerm, voteArgs.Id, voteArgs.LastSlotIndex, voteArgs.LastSlotTerm)

	if lastSlotTerm > voteArgs.LastSlotTerm ||
		(lastSlotTerm == voteArgs.LastSlotTerm && lastSlotIndex > int(voteArgs.LastSlotIndex)) {
		log.Printf("%d last slot CurrentTerm  is %d, last index is %d "+
			"later than vote's CurrentTerm %d, index %d reject vote",
			rf.me, lastSlotTerm, lastSlotIndex, voteArgs.LastSlotTerm, voteArgs.LastSlotIndex)
		reply.Result = voteReplyLatestLogEntryIsNotUpdateToMe
		return
	}

	if _, ok := rf.status.VoteFor[voteArgs.Term]; ok {
		log.Printf("%d CurrentTerm %d is already vote", rf.me, voteArgs.Term)
		reply.Result = voteReplyApplyAlreadyVote
		return
	}

	log.Printf("%d vote for service %d", rf.me, voteArgs.Id)

	rf.status.VoteFor[voteArgs.Term] = voteArgs.Id
	rf.persist()
	reply.Result = voteReplySuccess
	return
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
