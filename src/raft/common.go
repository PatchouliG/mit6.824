package raft

import (
	"log"
	"time"
)

//return Log index  if match
func (s *Status) logContain(index Index, term Term) bool {
	return len(s.Log) > int(index) && s.Log[index].Term == term
}

// not need any more, just compare current
//// -1 less 0 equal 1 large
//func (s *Status) lastLogTermCompare(CurrentTerm CurrentTerm) int {
//	if len(s.Log) == 0 {
//		return -1
//	}
//	entry := s.Log[len(s.Log)]
//	return int(entry.CurrentTerm - CurrentTerm)
//}

// return true if append success
func (rf *Raft) handleAppend(appendArg *AppendEntryArgs) AppendEntryReply {

	if rf.status.CurrentTerm > appendArg.CurrentTerm {
		log.Printf("%d append entry's CurrentTerm mistach, current CurrentTerm is %d, "+
			"CurrentTerm in requeset is %d ,reject it", rf.me, rf.status.CurrentTerm, appendArg.CurrentTerm)
		return AppendEntryReply{appendEntryStaleTerm}
	}

	// update current CurrentTerm
	log.Printf("%d receive append request ,update CurrentTerm from %d to %d", rf.me, appendArg.CurrentTerm,
		appendArg.CurrentTerm)
	rf.status.CurrentTerm = appendArg.CurrentTerm

	//update append time
	rf.lastAppendEntryTime = time.Now()
	match := rf.status.logContain(appendArg.PreviousEntryIndex, appendArg.PreviousEntryTerm)
	if match {
		if len(rf.status.Log)-1 != int(appendArg.PreviousEntryIndex) {
			log.Printf("delete all Log after the match")
			rf.status.Log = rf.status.Log[:appendArg.PreviousEntryIndex+1]
		}
		rf.status.Log = append(rf.status.Log, appendArg.Entries...)
		rf.persist()
		log.Printf("append finish, log size is %d", len(rf.status.Log))
		return AppendEntryReply{appendEntryAccept}
	}
	log.Printf("append entry prev entry match fail")
	return AppendEntryReply{appendEntryNotMatch}

}

func (rf *Raft) handleVote(voteArgs RequestVoteArgs) (reply RequestVoteReply) {

	log.Printf("%d receive vote", rf.me)
	if voteArgs.Term <= rf.status.CurrentTerm {
		log.Printf("%d receive request from an old CurrentTerm, current term is %d, reply term is %d, refuse it",
			rf.me, rf.status.CurrentTerm, voteArgs.Term)
		reply.Ok = false
		return
	}

	if _, ok := rf.status.VoteFor[voteArgs.Term]; ok {
		log.Printf("%d CurrentTerm %d is already vote", rf.me, voteArgs.Term)
		reply.Ok = false
		return
	}

	// update current
	log.Printf("%d update current CurrentTerm from %d to %d for vote", rf.me, rf.status.CurrentTerm, voteArgs.Term)
	rf.status.CurrentTerm = voteArgs.Term

	log.Printf("%d vote for service %d", rf.me, voteArgs.Id)

	rf.status.VoteFor[voteArgs.Term] = voteArgs.Id
	rf.persist()
	reply.Ok = true
	return
}
