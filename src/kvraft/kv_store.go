package kvraft

type KvStore struct {
	kvMap          map[string]string
	lastClientOpId map[int32]int
}
type KvStoreGetResult struct {
	value string
	exits bool
}

func newKVStore() KvStore {
	return KvStore{make(map[string]string), make(map[int32]int)}
}

func (kvs *KvStore) Get(key string) KvStoreGetResult {
	a, b := kvs.kvMap[key]
	return KvStoreGetResult{a, b}
}

func (kvs *KvStore) Put(opId OpId, key string, value string) {
	if kvs.checkAndSetOpId(opId) {
		kvs.kvMap[key] = value
	}
}
func (kvs *KvStore) append(opId OpId, key string, value string) {
	if kvs.checkAndSetOpId(opId) {
		kvs.kvMap[key] = kvs.kvMap[key] + value
	}
}

func (kvs *KvStore) checkAndSetOpId(opId OpId) bool {
	if value, ok := kvs.lastClientOpId[opId.ClientId]; ok {
		if value >= opId.Id {
			return false
		}
	}
	kvs.lastClientOpId[opId.ClientId] = opId.Id
	return true
}
