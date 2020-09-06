package raft

import (
	"log"
	"time"
)

func (rf *Raft) candidateRoutine() string {

	rf.status.CurrentTerm++
	rf.status.VoteFor[rf.status.CurrentTerm] = rf.me
	rf.persist()

	voteSuccesfulChan := make(chan int, 1000)
	voteSuccessfulCount := 0
	rf.sendVoteRequest(voteSuccesfulChan)

	for {
		select {
		case appendEntryArgs := <-rf.appendEntryRequest:
			res := rf.handleAppend(&appendEntryArgs)
			if res.Result != appendEntryStaleTerm {
				log.Printf("candidate recive a append entry request,change to follower")
				return follower
			}
		case voteArgs := <-rf.voteRequestChan:
			voteReply := rf.handleVote(voteArgs)
			rf.voteReplyChan <- voteReply
			if voteReply.Ok {
				return follower
			}
		case _ = <-voteSuccesfulChan:
			voteSuccessfulCount++
			// +1:plus self
			if voteSuccessfulCount+1 == len(rf.peers)/2+1 {
				log.Printf("%d get vote number %d,change to leader", rf.me, voteSuccessfulCount)
				return leader
			}

		case <-time.After(rf.randomElectionTimeout()):
			rf.candidateHandelTimeout(voteSuccesfulChan)
			voteSuccesfulChan = make(chan int)
			voteSuccessfulCount = 0
		}
	}
}

func (rf *Raft) candidateHandelTimeout(successChan chan int) {
	rf.status.CurrentTerm++
	rf.status.VoteFor[rf.status.CurrentTerm] = rf.me
	rf.persist()
	rf.sendVoteRequest(successChan)

}

func (rf *Raft) sendVoteRequest(successChan chan int) {
	log.Printf("%d candidate time out, begin CurrentTerm %d election", rf.me, rf.status.CurrentTerm)
	voteArgs := RequestVoteArgs{rf.status.CurrentTerm, Index(len(rf.status.Log) - 1), rf.me}
	for id, _ := range rf.peers {
		if id == rf.me {
			continue
		}
		go func(id int) {
			voteReply := RequestVoteReply{}
			log.Printf("%d send vote request from to %d", rf.me, id)
			rf.sendRequestVote(id, &voteArgs, &voteReply)
			if voteReply.Ok {
				log.Printf("%d receive vote success from %d", rf.me, id)
				successChan <- 0
			}
		}(id)
	}
}
