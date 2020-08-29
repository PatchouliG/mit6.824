package mr

import "time"

type workerInfo struct {
	lastHeatTime time.Time
	working      bool
	id           int64
}

func newWorkerInfo(id int64) workerInfo {
	return workerInfo{time.Now(), false, id}
}
