package mr

import (
	"fmt"
	"log"
	"time"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type Master struct {
	// Your definitions here.

}

// Your code here -- RPC handlers for the workerInfo to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func getAvailableWork() (int64, error) {
	return 0, fmt.Errorf("error")
}

const updateCilentInteval = clientKickOutTime / 5
const clientKickOutTime time.Duration = time.Second * 10

type workerInfo struct {
	lastHeatTime time.Time
	id           int64
}

var clientStatus = make(map[int64]workerInfo)

var heatBeatChan = make(chan int64, 1000)

func (m *Master) HeatBeat(args *int64, reply *int) error {
	*reply = 1
	//log.Printf("get heat beat from %d", *args)
	heatBeatChan <- *args
	return nil
}

//
// start a thread that listens for RPCs from workerInfo.go
//
func mainLoop() {
	t := time.NewTimer(updateCilentInteval)

	for {
		select {
		case workId := <-heatBeatChan:
			worker, ok := clientStatus[workId]
			if ok {
				worker.lastHeatTime = time.Now()
			} else {
				worker = workerInfo{id: workId, lastHeatTime: time.Now()}
			}
			clientStatus[workId] = worker

		//	check heathbeat
		case <-t.C:
			for key, v := range clientStatus {
				if time.Now().Sub(v.lastHeatTime) > clientKickOutTime {
					delete(clientStatus, key)
					log.Printf("time expire, kick out %d", key)
				}
			}
			t = time.NewTimer(updateCilentInteval)
		}
	}
}
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	go mainLoop()
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.

	m.server()
	return &m
}
