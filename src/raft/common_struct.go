package raft

type Index int
type Term int

type FollowerInfo struct {
	NextIndex  Index
	MatchIndex Index
}
type Status struct {
	CurrentTerm Term
	Log         []Entry
	VoteFor     map[Term]int
	NextIndex   Index
}

func (s *Status) getTerm(index Index) Term {
	return s.Log[index].Term
}

func NewStatus() Status {
	res := Status{0, []Entry{}, make(map[Term]int), 1}
	res.Log = append(res.Log, newEmptyEntry(0))
	return res
}

type Entry struct {
	Term       Term
	IsHeatBeat bool
	Command    Command
	Index      Index
}

func newEmptyEntry(term Term) Entry {
	return Entry{term, true, Command{}, -1}
}

//todo empty Command
type Command struct {
	Content interface{}
}

type PeerStatus struct {
	nextIndex  Index
	matchIndex Index
}

func (rf *Raft) newPeerStatus() PeerStatus {
	return PeerStatus{nextIndex: Index(len(rf.status.Log)), matchIndex: 0}
}

type StartReply struct {
	Index int
	Term  int
}
