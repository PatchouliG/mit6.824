package raft

import (
	"log"
	"time"
)

func (rf *Raft) followerRoutine() string {
	//rf.lastAppendEntryTime = time.Now()
	timeout := rf.randomElectionTimeout()
	timer := time.NewTimer(timeout)
	receiveAppend := false
	grantingVote := false
	for {
		select {
		case appendEntryArgs := <-rf.appendEntryRequest:
			if appendEntryArgs.CurrentTerm >= rf.status.CurrentTerm {
				receiveAppend = true
			}

			if appendEntryArgs.CurrentTerm > rf.status.CurrentTerm {
				// update current CurrentTerm
				log.Printf("%d receive append request ,update CurrentTerm from %d to %d", rf.me, rf.status.CurrentTerm,
					appendEntryArgs.CurrentTerm)
				rf.status.CurrentTerm = appendEntryArgs.CurrentTerm
				rf.persist()
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
			timeout = rf.randomElectionTimeout()
			timer = time.NewTimer(timeout)
			receiveAppend = false
			grantingVote = false
		}
	}
}
