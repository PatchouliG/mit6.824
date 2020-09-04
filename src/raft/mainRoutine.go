package raft

func (rf *Raft) mainRoutine() {
	rf.callRoleRoutine(follower)
}

func (rf *Raft) callRoleRoutine(role string) {
	for {
		switch role {
		case follower:
			role = rf.followerRoutine()
		case candidate:
			role = rf.candidateRoutine()
		case leader:
			role = rf.leaderRoutine()
		}
	}
}
