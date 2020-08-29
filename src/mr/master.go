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

	files  []string
	jsList []jobStatus
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
			err := m.jsList[i].updateStatus(running)
			if err == nil {
				*reply = m.jsList[i].info
				return nil
			}
		}
	}
	return fmt.Errorf("no more work")
}

func (m *Master) jobDone(arg *JobInfo, reply *int) error {
	*reply = 0
	id := arg.id()
	for _, js := range m.jsList {
		if js.info.id() == id {
			err := js.updateStatus(finish)

			if err == nil {
				if arg.JobType == mapJob {
					m.jsList = append(m.jsList, newJobStatus(mapOutputFile(arg.InputFile[0], arg.ReduceNumber),
						arg.ReduceNumber, reduceJob))
					return nil
				}
			}
			log.Fatalf("update job %s status to finish error", js.info.id())
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
			for _, js := range m.jsList {
				if js.status == running && time.Now().Sub(js.startTime) > clientKickOutTime {
					js.updateStatus(pending)
					log.Printf("time expire, set job to expire")
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
	for _, i := range files {
		jsList = append(jsList, newJobStatus([]string{i}, nReduce, mapJob))
	}
	fmt.Printf("%+v", jsList)
	m := Master{files, jsList}

	// Your code here.

	m.server()
	return &m
}
