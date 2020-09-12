package raft

import (
	"log"
	"time"
)

func (rf *Raft) followerRoutine() string {
	timeout := rf.randomElectionTimeout()
	timer := time.NewTimer(timeout)
	receiveAppend := false
	grantingVote := false
	for {
		select {
		case appendEntryArgs := <-rf.appendEntryRequest:
			if appendEntryArgs.CurrentTerm >= rf.status.CurrentTerm {
				receiveAppend = true
				log.Printf("%d recevie append,accept it", rf.me)
			}

			if appendEntryArgs.CurrentTerm > rf.status.CurrentTerm {
				// update current CurrentTerm
				//log.Printf("%d receive append request ,update CurrentTerm from %d to %d", rf.me, rf.status.CurrentTerm,
				//	appendEntryArgs.CurrentTerm)
				rf.setCurrentTerm(appendEntryArgs.CurrentTerm)
				rf.persist()
			}

			res := rf.handleAppend(&appendEntryArgs)
			rf.appendEntryReply <- res
		case voteArgs := <-rf.voteRequestChan:
			res := rf.handleVote(voteArgs)
			rf.voteReplyChan <- res

			if voteArgs.Term > rf.status.CurrentTerm {
				rf.setCurrentTerm(voteArgs.Term)
				rf.persist()
			}
			if res.Result == voteReplySuccess {
				grantingVote = true
				rf.persist()
			}

		case <-timer.C:
			if rf.killed() {
				log.Printf("%d is killed,exit", rf.me)
				return "exit"
			}
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
