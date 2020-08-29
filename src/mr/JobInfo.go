package mr

import "strconv"

const (
	mapJob    = "map"
	reduceJob = "reduce"
)

type JobInfo struct {
	JobType      string
	InputFile    []string
	ReduceNumber int
}

//func (ji *JobInfo) tmpOutputFile(file string) string {
//	return "tmp_" + ji.outputFile(file)
//}
func mapOutputFile(file string, reduceNumber int) []string {
	var res []string
	for i := 0; i < reduceNumber; i++ {
		res = append(res, "mr-"+file+"-"+strconv.Itoa(i))
	}
	return res
}

func reduceOutputFile(reduceNumber int) string {
	return "mr-out-" + strconv.Itoa(reduceNumber)
}

func (ji *JobInfo) id() string {
	res := ""
	for _, file := range ji.InputFile {
		res += file
		res += " "

	}
	return res

}
