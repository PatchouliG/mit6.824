package mr

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func marshalToFile(absFile string, objects []KeyValue) {
	_,fileName := filepath.Split(absFile)
	cw, _ := os.Getwd()

	file, error := os.Create(path.Join(cw, fileName))
	defer file.Close()
	if error != nil {
		log.Panicf("open file fail %s", error)
	}
	for _, object := range objects {
		result, error := json.Marshal(object)
		if error != nil {
			log.Panicf("march fail")
		}
		file.Write(result)
		file.WriteString("\n")
	}
}

func unmarshalFromFile(fileName string) []KeyValue {
	cw, _ := os.Getwd()
	file, error := os.Open(path.Join(cw, fileName))
	defer file.Close()
	if error != nil {
		log.Panicf("open file fail")
	}
	var res []KeyValue
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var tmp KeyValue
		line := scanner.Bytes()
		json.Unmarshal(line, &tmp)
		if scanner.Err() != nil {
			log.Fatalf(" > Failed with error %v\n", scanner.Err())
		}
		res = append(res, tmp)
	}
	return res
}

func findReduceInputFile(dir string, id int) []string {
	var res []string
	files, err := ioutil.ReadDir("/Users/wn/code/6.824-golabs-2020/src/main")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), "txt-"+strconv.Itoa(id)) {
			res = append(res, f.Name())
		}
	}
	return res
}
