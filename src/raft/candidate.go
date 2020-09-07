package raft

import (
	"log"
	"time"
)

func (rf *Raft) candidateRoutine() string {

	rf.status.CurrentTerm++
	rf.status.VoteFor[rf.status.CurrentTerm] = rf.me
	rf.persist()

	//log.Printf("%d debug inc term %d", rf.me, rf.status.CurrentTerm)

	voteSuccesfulChan := make(chan RequestVoteReply, 1000)
	voteSuccessfulCount := 0
	rf.sendVoteRequest(voteSuccesfulChan)
	timer := time.NewTimer(rf.randomElectionTimeout())

	for {
		select {
		case appendEntryArgs := <-rf.appendEntryRequest:
			res := rf.handleAppend(&appendEntryArgs)
			rf.appendEntryReply <- res
			if appendEntryArgs.CurrentTerm > rf.status.CurrentTerm {
				log.Printf("%d as candidate recive a append request, it's term is later, turn to follower", rf.me)
				rf.status.CurrentTerm = appendEntryArgs.CurrentTerm
				return follower
			}
		case voteArgs := <-rf.voteRequestChan:
			voteReply := rf.handleVote(voteArgs)
			rf.voteReplyChan <- voteReply
			if voteArgs.Term > rf.status.CurrentTerm {
				rf.status.CurrentTerm = voteArgs.Term
				return follower
			}
		case voteReply := <-voteSuccesfulChan:
			switch voteReply.Result {
			case voteReplySuccess:
				voteSuccessfulCount++
				// +1:plus self
				if voteSuccessfulCount+1 == len(rf.peers)/2+1 {
					log.Printf("%d get vote number %d,change to leader", rf.me, voteSuccessfulCount)
					return leader
				}
			case voteReplyApplyAlreadyVote:
				continue
			case voteReplyLatestLogEntryIsNotUpdateToMe:
				continue
			case voteReplyStaleTerm:
				log.Printf("%d receive stale CurrentTerm when request vote, turn to follower", rf.me)
				rf.status.CurrentTerm = voteReply.CurrentTerm
				return follower
			}

		case <-timer.C:
			log.Printf("%d canditate time out ", rf.me)
			rf.candidateHandelTimeout(voteSuccesfulChan)
			voteSuccesfulChan = make(chan RequestVoteReply, 1000)
			voteSuccessfulCount = 0
			timer = time.NewTimer(rf.randomElectionTimeout())
		}
	}
}

func (rf *Raft) candidateHandelTimeout(successChan chan RequestVoteReply) {
	rf.status.CurrentTerm++
	rf.status.VoteFor[rf.status.CurrentTerm] = rf.me
	rf.persist()
	rf.sendVoteRequest(successChan)

}

func (rf *Raft) sendVoteRequest(successChan chan RequestVoteReply) {
	voteArgs := RequestVoteArgs{rf.status.CurrentTerm,
		Index(len(rf.status.Log) - 1),
		rf.status.Log[len(rf.status.Log)-1].Term,
		rf.me}
	for id, _ := range rf.peers {
		if id == rf.me {
			continue
		}
		go func(id int) {
			voteReply := RequestVoteReply{}
			log.Printf("%d send vote request from to %d", rf.me, id)
			rf.sendRequestVote(id, &voteArgs, &voteReply)
			successChan <- voteReply
		}(id)
	}
	//log.Printf("debug %d send vote end,current term is %d", rf.me, rf.status.CurrentTerm)
}
