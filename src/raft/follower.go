package raft

import (
	"log"
	"time"
)

func (rf *Raft) followerRoutine() string {
	for {
		select {
		case appendEntryArgs := <-rf.appendEntryIn:
			res := rf.handleAppend(&appendEntryArgs)
			rf.appendEntryout <- AppendEntryReply{res}
		case voteArgs := <-rf.voteArgs:
			rf.handleVote(voteArgs)
		case <-time.After(electionTimeout):
			role := rf.followerHandleTimeout()
			if role != follower {
				return role
			}
		}
	}
}

//func (rf *Raft) handleAppend(appendArg *AppendEntryArgs) bool {
//
//	rf.mu.Lock()
//	defer rf.mu.Unlock()
//
//	if rf.status.currentTerm > appendArg.Term {
//		log.Printf("append entry's Term mistach, current Term is %d, "+
//			"Term in requeset is %d ,reject it", rf.status.currentTerm, appendArg.Term)
//		return false
//	}
//	//update append time
//	rf.lastAppendEntryTime = time.Now()
//	index, found := rf.status.logMatchAt(appendArg.PreviousEntryIndex)
//	if found {
//		if len(rf.status.log)-1 != index {
//			log.Printf("delete all log after the match")
//			rf.status.log = rf.status.log[:index+1]
//		}
//		rf.status.log = append(rf.status.log, appendArg.Entries...)
//		rf.persist()
//		return true
//	}
//	log.Printf("append entry prev entry match fail")
//	return false
//
//}

// return role
func (rf *Raft) followerHandleTimeout() string {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if time.Now().Sub(rf.lastAppendEntryTime) > electionTimeout {
		log.Printf("timeout, begin election")
		return candidate
	}
	return follower
}
