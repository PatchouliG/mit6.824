package raft

import (
	"log"
	"time"
)

func (rf *Raft) followerRoutine() string {
	for {
		select {
		case appendEntryArgs := <-rf.appendEntryRequest:
			res := rf.handleAppend(&appendEntryArgs)
			rf.appendEntryReply <- res
		case voteArgs := <-rf.voteRequestChan:
			res := rf.handleVote(voteArgs)
			rf.voteReplyChan <- res
		case <-time.After(timeoutCheck):
			nextRole := rf.followerCheckTimeout()
			if nextRole {
				return candidate
			}
		}
	}
}

// true if change
func (rf *Raft) followerCheckTimeout() bool {

	if time.Now().Sub(rf.lastAppendEntryTime) > rf.randomElectionTimeout() {
		log.Printf("%d timeout, begin election,current CurrentTerm is %d", rf.me, rf.status.CurrentTerm)
		return true
	}
	return false
}
