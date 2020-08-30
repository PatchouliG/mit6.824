package mr

import (
	"log"
	"sync"
	"time"
)

const un_allocate = -1
const workDone = -2

const (
	pending = 1
	running = 2
	finish  = 3
)

type jobStatus struct {
	info      JobInfo
	status    int
	startTime time.Time
	mux       sync.Mutex
}

func newMapJobStatus(files string, reduceNumber int) jobStatus {
	jobinfo := newMapJobInfo(files, reduceNumber)
	return jobStatus{jobinfo, pending, time.Now(), sync.Mutex{}}
}
func newReduceJobStatus(reduceNumber int) jobStatus {
	jobinfo := newReduceJobInfo(reduceNumber)
	return jobStatus{jobinfo, pending, time.Now(), sync.Mutex{}}
}

func (js *jobStatus) updateStatus(status int) {
	js.mux.Lock()
	defer js.mux.Unlock()

	if status-js.status == 1 || (js.status == running && status == pending) {
		js.status = status
		if status == running {
			js.startTime = time.Now()
		}
	} else {
		log.Fatalf("update %s status error, from %d to %d", js.info.id(), js.status, status)
	}
}
