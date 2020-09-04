package raft

import (
	"log"
	"time"
)

func (rf *Raft) candidateRoutine() string {
	for {
		rf.mu.Lock()
		rf.status.currentTerm++
		rf.status.voteFor[rf.status.currentTerm] = rf.me
		rf.persist()
		rf.mu.Unlock()
		select {
		case appendEntryArgs := <-rf.appendEntryIn:
			res := rf.handleAppend(&appendEntryArgs)
			if res {
				log.Printf("candidate recive a append entry request,change to follower")
				return follower
			}
		case voteArgs := <-rf.voteArgs:
			rf.handleVote(voteArgs)
		case <-time.After(electionTimeout):
			rf.candidateHandelTimeout()
		}
	}
}

func (rf *Raft) candidateHandelTimeout() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.status.currentTerm++
	rf.persist()

}
