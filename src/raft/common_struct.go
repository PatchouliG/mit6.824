package raft

type Index int
type Term int

type FollowerInfo struct {
	nextIndex  Index
	matchIndex Index
}
type Status struct {
	CurrentTerm Term
	Log         []Entry
	VoteFor     map[Term]int
}

func (s *Status) getTerm(index Index) Term {
	return s.Log[index].Term
}

func NewStatus() Status {
	res := Status{0, []Entry{}, make(map[Term]int)}
	res.Log = append(res.Log, newEmptyEntry(0))
	return res
}

type Entry struct {
	Term    Term
	Command Command
}

func newEmptyEntry(term Term) Entry {
	return Entry{term, Command{}}
}

//todo empty Command
type Command struct {
}

type PeerStatus struct {
	nextIndex  Index
	matchIndex Index
}

func (rf *Raft) newPeerStatus() PeerStatus {
	return PeerStatus{nextIndex: Index(len(rf.status.Log)), matchIndex: 0}
}
