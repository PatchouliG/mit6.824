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
}

//return Log index  if match
func (s *Status) logContain(index Index, term Term) bool {
	return len(s.Log) > int(index) && s.Log[index].Term == term
}

func (s *Status) lastIndex() Index {
	return s.Log[len(s.Log)-1].Index
}
func (s *Status) lastEntry() Entry {
	return s.Log[len(s.Log)-1]
}

func (s *Status) lastEntryContainsOperation() (Entry, bool) {
	for i := len(s.Log) - 1; i >= 0; i-- {
		if !s.Log[i].IsHeatBeat {
			return s.Log[i], true
		}
	}
	return Entry{}, false
}

func (s *Status) getEntry(index Index) (Entry, int, bool) {
	if len(s.Log) == 0 {
		return Entry{}, 0, false
	}
	position := int(index) - int(s.Log[0].Index)
	if position >= len(s.Log) {
		return Entry{}, 0, false
	}
	return s.Log[position], position, true
}

func (s *Status) deleteAfter(position int) {
	s.Log = s.Log[:position]
}
func (s *Status) append(entries []Entry) {
	s.Log = append(s.Log, entries...)
}

func (s *Status) getTerm(index Index) Term {
	return s.Log[index].Term
}

//func (s *Status) lastEntry() Entry {
//	return s.Log[len(s.Log)-1]
//}

func NewStatus() Status {
	res := Status{0, []Entry{}, make(map[Term]int)}
	res.Log = append(res.Log, newEmptyEntry(0, 0))
	return res
}

type Entry struct {
	Term       Term
	IsHeatBeat bool
	Command    Command
	Index      Index
}

func newEmptyEntry(term Term, index Index) Entry {
	return Entry{term, true, Command{}, index}
}

//todo empty Command
type Command struct {
	Operation interface{}
	//the tester need the command Index
	Index int
}

type PeerStatus struct {
	nextIndex  Index
	matchIndex Index
}

type StartReply struct {
	Index int
	Term  int
}
