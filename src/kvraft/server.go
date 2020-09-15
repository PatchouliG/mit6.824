package kvraft

import (
	"../labgob"
	"../labrpc"
	"../raft"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

type Op struct {
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.

	Key    string
	Value  string
	OpId   OpId
	OpType string
}

const commandTimeout = time.Duration(time.Second)

type KVServer struct {
	mu      sync.Mutex
	me      int
	rf      *raft.Raft
	applyCh chan raft.ApplyMsg
	dead    int32 // set by Kill()

	maxraftstate int // snapshot if log grows this big

	// Your definitions here.

	kvs KvStore
	// OpId->value
	replyNotify   map[OpId]chan KvStoreGetResult
	replyNotifyMu sync.Mutex
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	c := kv.registerResultNotify(args.OpId)
	defer kv.unRegisterResultNotify(args.OpId)

	_, _, isLeader := kv.rf.Start(Op{args.Key, "", args.OpId, "Get"})
	if !isLeader {
		reply.Err = ErrWrongLeader
		return
	}
	timer := time.NewTimer(commandTimeout)
	select {
	case result := <-c:
		value, ok := result.value, result.exits
		if !ok {
			reply.Err = ErrNoKey
		} else {
			reply.Err = OK
			reply.Value = value
		}
	case <-timer.C:
		reply.Err = ErrWrongLeader
	}
}

func (kv *KVServer) PutAppend(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	c := kv.registerResultNotify(args.OpId)
	defer kv.unRegisterResultNotify(args.OpId)

	_, _, isLeader := kv.rf.Start(Op{args.Key, args.Value, args.OpId, args.Op})
	if !isLeader {
		reply.Err = ErrWrongLeader
		return
	}
	timer := time.NewTimer(commandTimeout)
	select {
	case result := <-c:
		_, ok := result.value, result.exits
		if !ok {
			reply.Err = ErrNoKey
		} else {
			reply.Err = OK
		}
	case <-timer.C:
		reply.Err = ErrWrongLeader
	}
}

//func (kv *KVServer) registerResp(ClientId int, opId int) chan struct{} {
//	c := make(chan struct{})
//	kv.replyDone[strconv.Itoa(ClientId)+":"+strconv.Itoa(opId)] = c
//	return c
//}
//func (kv *KVServer) unRegisterResp(ClientId int, opId int) {
//	kv.replyDone[strconv.Itoa(ClientId)+":"+strconv.Itoa(opId)] = nil
//}

//
// the tester calls Kill() when a KVServer instance won't
// be needed again. for your convenience, we supply
// code to set rf.dead (without needing a lock),
// and a killed() method to test rf.dead in
// long-running loops. you can also add your own
// code to Kill(). you're not required to do anything
// about this, but it may be convenient (for example)
// to suppress debug output from a Kill()ed instance.
//
func (kv *KVServer) Kill() {
	atomic.StoreInt32(&kv.dead, 1)
	kv.rf.Kill()
	// Your code here, if desired.
}

func (kv *KVServer) killed() bool {
	z := atomic.LoadInt32(&kv.dead)
	return z == 1
}

func (kv *KVServer) notifyResult(opId OpId, result KvStoreGetResult) {
	kv.replyNotifyMu.Lock()
	defer kv.replyNotifyMu.Unlock()
	if replyChan, ok := kv.replyNotify[opId]; ok {
		replyChan <- result
	}
}
func (kv *KVServer) registerResultNotify(opId OpId) chan KvStoreGetResult {
	res := make(chan KvStoreGetResult, 2)
	kv.replyNotifyMu.Lock()
	defer kv.replyNotifyMu.Unlock()
	kv.replyNotify[opId] = res
	return res
}
func (kv *KVServer) unRegisterResultNotify(opId OpId) {
	kv.replyNotifyMu.Lock()
	defer kv.replyNotifyMu.Unlock()
	delete(kv.replyNotify, opId)
}

func (kv *KVServer) applyResultRoutine() {
	for {
		timer := time.NewTimer(time.Second)
		select {
		case applyMsg := <-kv.applyCh:
			op := applyMsg.Command.(Op)
			kv.replyNotifyMu.Lock()
			notifyChan, ok := kv.replyNotify[op.OpId]
			kv.replyNotifyMu.Unlock()
			switch op.OpType {
			case OpGet:
				res := kv.kvs.Get(op.Key)

				if ok {
					notifyChan <- res
				}
			case OpPut:
				kv.kvs.Put(op.OpId, op.Key, op.Value)
				if ok {
					notifyChan <- KvStoreGetResult{"", true}
				}
			case OpAppend:
				kv.kvs.append(op.OpId, op.Key, op.Value)
				if ok {
					notifyChan <- KvStoreGetResult{"", true}
				}
			}
		case <-timer.C:
			timer = time.NewTimer(time.Second)
			if kv.killed() {
				return
			}
		}

	}
}

//
// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots through the underlying Raft
// implementation, which should call persister.SaveStateAndSnapshot() to
// atomically save the Raft state along with the snapshot.
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
//
func StartKVServer(servers []*labrpc.ClientEnd, me int, persister *raft.Persister, maxraftstate int) *KVServer {
	// call labgob.Register on structures you want
	// Go's RPC library to marshall/unmarshall.
	labgob.Register(Op{})

	kv := new(KVServer)
	kv.me = me
	kv.maxraftstate = maxraftstate
	kv.replyNotify = make(map[OpId]chan KvStoreGetResult)

	// You may need initialization code here.

	kv.applyCh = make(chan raft.ApplyMsg)
	kv.rf = raft.Make(servers, me, persister, kv.applyCh)
	//kv.replyDone = make(map[string]chan struct{})
	kv.kvs = newKVStore()
	go kv.applyResultRoutine()

	// You may need initialization code here.
	//go kv.dispatchResult()

	return kv
}
