package kvraft

const (
	OK             = "OK"
	ErrNoKey       = "ErrNoKey"
	ErrWrongLeader = "ErrWrongLeader"
)

const (
	OpPut    = "Put"
	OpGet    = "Get"
	OpAppend = "Append"
)

type Err string

type OpId struct {
	ClientId int32
	Id       int
}

// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	Op    string // "Put" or "Append"
	// You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	OpId OpId
}

type PutAppendReply struct {
	Err Err
}

type GetArgs struct {
	OpId OpId
	Key  string
	// You'll have to add definitions here.
}

type GetReply struct {
	Err   Err
	Value string
}
