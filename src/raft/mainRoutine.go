package raft

import "log"

func (rf *Raft) mainRoutine() {
	log.Printf("%d begin main routine", rf.me)
	rf.callRoleRoutine(follower)
}

func (rf *Raft) callRoleRoutine(role string) {
	for {
		log.Printf("%d run as %s", rf.me, role)
		switch role {
		case follower:
			role = rf.followerRoutine()
		case candidate:
			role = rf.candidateRoutine()
		case leader:
			rf.isLeader = true
			role = rf.leaderRoutine()
			rf.isLeader = false
		default:
			return
		}
	}
}
