package mr

import (
	"fmt"
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

func newJobStatus(files []string, reduceNumber int, jobType string) jobStatus {
	return jobStatus{JobInfo{
		jobType, files, reduceNumber,
	}, pending, time.Now(), sync.Mutex{}}
}

func (js *jobStatus) updateStatus(status int) error {
	js.mux.Lock()
	defer js.mux.Unlock()

	if status-js.status == 1 || (js.status == running && status == pending) {
		js.status = status
		return nil
	}
	return fmt.Errorf("status update fail, from %d to %d",
		status, js.status)
}
