package raft

import (
	"log"
	"time"
)

//return log index  if match
func (s *Status) logMatchAt(index Index, term Term) bool {
	return len(s.log) > int(index) && s.log[index].term == term
}

func (rf *Raft) handleAppend(appendArg *AppendEntryArgs) bool {
	defer rf.mu.Unlock()

	if rf.status.currentTerm > appendArg.Term {
		log.Printf("append entry's Term mistach, current Term is %d, "+
			"Term in requeset is %d ,reject it", rf.status.currentTerm, appendArg.Term)
		return false
	}
	//update append time
	rf.lastAppendEntryTime = time.Now()
	match := rf.status.logMatchAt(appendArg.PreviousEntryIndex, appendArg.PreviousEntryTerm)
	if match {
		if len(rf.status.log)-1 != int(appendArg.PreviousEntryIndex) {
			log.Printf("delete all log after the match")
			rf.status.log = rf.status.log[:appendArg.PreviousEntryIndex+1]
		}
		rf.status.log = append(rf.status.log, appendArg.Entries...)
		rf.persist()
		return true
	}
	log.Printf("append entry prev entry match fail")
	return false

}

func (rf *Raft) handleVote(voteArgs RequestVoteArgs) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	log.Print("receive vote")
	reply := RequestVoteReply{}
	if voteArgs.term <= rf.status.currentTerm {
		log.Print("receive request from an old Term, refuse it")
		reply.ok = false
		//return
	}
	if _, ok := rf.status.voteFor[voteArgs.term]; ok {
		log.Printf("Term %d is already vote", voteArgs.term)
		reply.ok = false
		//return
	}
	rf.status.voteFor[voteArgs.term] = voteArgs.id
	rf.persist()
	rf.voteReply <- reply
}

func (rf *Raft) getAppendArgForHeatBeat() AppendEntryArgs {
	//todo
	return AppendEntryArgs{}
}
