package mr

import (
	"path"
	"path/filepath"
	"strconv"
)

const (
	mapJob    = "map"
	reduceJob = "reduce"
)

type JobInfo struct {
	JobType      string
	InputFile    []string
	ReduceNumber int
	ReduceId     int
}

func newMapJobInfo(file string, reduceNumber int) JobInfo {
	return JobInfo{mapJob, []string{file}, reduceNumber, -1}

}

func newReduceJobInfo(reduceId int) JobInfo {
	return JobInfo{reduceJob, []string{}, -1, reduceId}
}

func mapOutputFile(file string, reduceNumber int) []string {
	var res []string
	dir, fileName := filepath.Split(file)
	for i := 0; i < reduceNumber; i++ {
		res = append(res, path.Join(dir, "mr-"+fileName+"-"+strconv.Itoa(i)))
	}
	return res
}

func reduceOutputFile(reduceNumber int) string {
	return "mr-out-" + strconv.Itoa(reduceNumber)
}

func (ji *JobInfo) id() string {
	if ji.JobType == mapJob {
		res := ""
		for _, file := range ji.InputFile {
			res += file
			res += " "
		}
		return res
	} else {
		return "reduce" + strconv.Itoa(ji.ReduceId)
	}
}
