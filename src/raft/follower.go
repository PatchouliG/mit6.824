package raft

import (
	"time"
)

func (rf *Raft) followerRoutine() string {
	//rf.lastAppendEntryTime = time.Now()
	for {
		receiveAppend := false
		grantingVote := false
		timeout := rf.randomElectionTimeout()
		timer := time.NewTimer(timeout)
		select {
		case appendEntryArgs := <-rf.appendEntryRequest:
			if appendEntryArgs.CurrentTerm >= rf.status.CurrentTerm {
				receiveAppend = true
				//rf.lastAppendEntryTime = time.Now()
			}
			res := rf.handleAppend(&appendEntryArgs)
			rf.appendEntryReply <- res
		case voteArgs := <-rf.voteRequestChan:
			res := rf.handleVote(voteArgs)
			rf.voteReplyChan <- res

			if voteArgs.Term > rf.status.CurrentTerm {
				rf.status.CurrentTerm = voteArgs.Term
				grantingVote = true
				rf.persist()
			}

		case <-timer.C:
			if !receiveAppend && !grantingVote {
				return candidate
			}
		}
	}
}
