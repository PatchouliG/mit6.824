package kvraft

import (
	"../labrpc"
	"crypto/rand"
	"math/big"
	"sync/atomic"
	"time"
)

var nextClerkId int32 = 0

type Clerk struct {
	servers []*labrpc.ClientEnd
	// You will have to modify this struct.
	leaderId int
	nextOpId int
	clerkId  int32
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	ck.leaderId = 0

	ck.clerkId = atomic.AddInt32(&nextClerkId, 1)
	ck.nextOpId = 0
	// You'll have to add code here.
	return ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) Get(key string) string {
	opId := OpId{ck.clerkId, ck.nextOpId}
	args := GetArgs{opId, key}
	reply := GetReply{}
	ck.nextOpId++
	for {
		ok := ck.servers[ck.leaderId].Call("KVServer.Get", &args, &reply)
		if ok {
			switch reply.Err {
			case OK:
				return reply.Value
			case ErrWrongLeader:
			case ErrNoKey:
				return ""
			}
		}
		ck.tryNextLeader()
		time.Sleep(time.Millisecond * 50)
	}
	// You will have to modify this function.
}

func (ck *Clerk) tryNextLeader() int {
	ck.leaderId = (ck.leaderId + 1) % len(ck.servers)
	return ck.leaderId
}

//
// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) PutAppend(key string, value string, op string) {
	opId := OpId{ck.clerkId, ck.nextOpId}
	args := PutAppendArgs{key, value, op, opId}
	reply := PutAppendReply{}
	ck.nextOpId++
	for {
		ok := ck.servers[ck.leaderId].Call("KVServer.PutAppend", &args, &reply)
		if ok {
			switch reply.Err {
			case OK:
				return
			case ErrNoKey:
				return
			case ErrWrongLeader:
			}
		}
		ck.tryNextLeader()
		time.Sleep(time.Millisecond * 50)
	}
	// You will have to modify this function.
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
