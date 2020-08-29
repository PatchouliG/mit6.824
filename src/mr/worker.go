package mr

import (
	"fmt"
	"time"
)
import "log"
import "net/rpc"
import "hash/fnv"

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
//
//func heatBeatCall(id int64) {
//
//	for {
//		var reply int
//		fmt.Printf("test")
//
//		send the RPC request, wait for the reply.
//call("Master.HeatBeat", &id, &reply)
//time.Sleep(time.Second * 1)
//}
//
//}
func callFetchJob() (JobInfo, error) {
	arg := 0
	reply := JobInfo{}
	res := call("Master.FetchJob", &arg, &reply)
	if !res {
		return reply, fmt.Errorf("fetch job fail")
	}
	return reply, nil
}

func callJobDone(info JobInfo) {

}
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your workerInfo implementation here.

	// uncomment to send the Example RPC to the master.
	//CallExample()
	for {
		ji, error := callFetchJob()
		if error != nil {
			log.Print(error)
			time.Sleep(time.Second * 3)
			continue
		}
		log.Printf("get file %s", ji.InputFile)
		if ji.JobType == mapJob {
			//doMapJob(ji, mapf)

		} else {
			//	todo
			//callJobDone()
		}
	}

}

//func doMapJob(ji JobInfo, mapf func(string, string) []KeyValue) {
//	filename := ji.InputFile[0]
//	file, err := os.Open(filename)
//	if err != nil {
//		log.Fatalf("cannot open %v", filename)
//	}
//	content, err := ioutil.ReadAll(file)
//	if err != nil {
//		log.Fatalf("cannot read %v", filename)
//	}
//	file.Close()
//	kva := mapf(filename, string(content))
//	var reduceBuffer [ji.ReduceNumber][]KeyValue
//	for _, kv := range kva {
//		number := ihash(kv.Key)
//		reduceBuffer[number] = append(reduceBuffer[number], kv)
//	}
//	outputFiles := mapOutputFile(ji.InputFile[0], ji.ReduceNumber)
//	for i, fileName := range outputFiles {
//		file, err := os.Create(fileName)
//		if err != nil {
//			log.Fatalf("cannot open %v", filename)
//		}
//		file.WriteString()
//
//	}
//
//	//todo
//	callJobDone()
//}

//
// example function to show how to make an RPC call to the master.
//
// the RPC argument and reply types are defined in rpc.go.
//

func fetchJob() {

}
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	call("Master.Example", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
