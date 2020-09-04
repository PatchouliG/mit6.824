package raft

import "time"

func (rf *Raft) leaderRoutine() string {
	args := rf.getAppendArgForHeatBeat()
	reply := AppendEntryReply{}
	for _, peer := range rf.peers {
		// todo go on here
		//todo handle reply async
		//go rf.sendAppendEntry(peer,[]Command{newEmptyEntry()})
	}
	//todo recored peer log status,like next entry
	for {
		select {
		case appendEntryArgs := <-rf.appendEntryIn:
			// todo change to follower
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
	return ""
}

// run by go routine
func (rf *Raft) sendAppendEntry(peerId int, entry Entry) {

}
