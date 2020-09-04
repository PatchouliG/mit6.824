package raft

type Index int
type Term int

type Status struct {
	currentTerm Term
	log         []Entry
	voteFor     map[Term]int
}

func NewStatus() Status {
	return Status{0, []Entry{}, make(map[Term]int)}
}

type Entry struct {
	term    Term
	command Command
}

func newEmptyEntry(term Term) Entry {
	return Entry{term, Command{}}
}

//todo empty command
type Command struct {
}

type PeerStatus struct {
	nextIndex  Index
	matchIndex Index
}

func (rf *Raft) newPeerStatus() PeerStatus {
	return PeerStatus{nextIndex: Index(len(rf.status.log)), matchIndex: 0}
}
