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

const updateCilentInteval = clientKickOutTime / 5
const clientKickOutTime time.Duration = time.Second * 10

type Master struct {
	// Your definitions here.

	files            []string
	jsList           []jobStatus
	finishedMapCount int
	mapJobNumber     int
	reduceNumber     int
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

// error if no available job found
func (m *Master) FetchJob(unUsed *int, reply *JobInfo) error {
	for i, js := range m.jsList {
		if js.status == pending {
			m.jsList[i].updateStatus(running)
			*reply = m.jsList[i].info
			log.Printf("send job %s", m.jsList[i].info.id())
			return nil
		}
	}
	return fmt.Errorf("no more work")
}

func (m *Master) JobDone(arg *JobInfo, reply *int) error {
	*reply = 0
	id := arg.id()
	log.Printf("rececive job done %s,job type %s", arg.id(), arg.JobType)
	for i, js := range m.jsList {
		if js.info.id() == id {
			m.jsList[i].updateStatus(finish)
			if arg.JobType == mapJob {
				m.finishedMapCount += 1
				if m.finishedMapCount < m.mapJobNumber {
				} else {
					log.Printf("all map job is finishd")
					for i := 0; i < m.reduceNumber; i++ {
						m.jsList = append(m.jsList, newReduceJobStatus(i))
					}
				}
			}
			return nil
		}
	}
	log.Fatalf("job id %s not found", arg.id())
	return nil
}

//func (m *Master) addReduceJob(file []string) {
//
//}

//func getReduceName(number int) string {
//	return "reduce_" + strconv.Itoa(number)
//}

func (m *Master) checkTimeOut() {
	t := time.NewTimer(updateCilentInteval)
	for {
		select {
		case <-t.C:
			for i, js := range m.jsList {
				if js.status == running && time.Now().Sub(js.startTime) > clientKickOutTime {
					m.jsList[i].updateStatus(pending)
					log.Printf("time expire, set job %s to expire", js.info.id())
				}
			}
			t = time.NewTimer(updateCilentInteval)
		}
	}
}

//
// start a thread that listens for RPCs from workerInfo.go
//

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
	go m.checkTimeOut()
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := true
	for _, js := range m.jsList {
		if js.status != finish {
			ret = false
			break
		}
	}

	// Your code here.

	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	//var reduceFils []string
	var jsList []jobStatus
	for _, file := range files {
		log.Printf(file)
		jsList = append(jsList, newMapJobStatus(file, nReduce))
	}
	a, _ := os.Getwd()
	fmt.Printf("debug cwd is %s", a)
	m := Master{files, jsList, 0, len(jsList), nReduce}

	// Your code here.

	m.server()
	return &m
}
