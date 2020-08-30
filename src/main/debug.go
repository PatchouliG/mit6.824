package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type test struct {
	A, B, C int
}

func marshalToFile(fileName string, objects []interface{}) {
	file, error := os.Create(fileName)
	defer file.Close()
	if error != nil {
		log.Panicf("open file fail")
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

func unmarshalFromFile(fileName string) []test {
	file, error := os.Open(fileName)
	defer file.Close()
	if error != nil {
		log.Panicf("open file fail")
	}
	var res []test
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var tmp test
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

func main() {
	//fmt.Println(findReduceInputFile("/Users/wn/code/6.824-golabs-2020/src/main", 3))

	//fltB, _ := json.Marshal(test{1, 2, 3})
	//file, error := os.Create("/Users/wn/code/6.824-golabs-2020/src/main/tmp")
	//if error != nil {
	//	log.Panic(error)
	//}
	//file.Write(fltB)
	//file.WriteString("\n")
	//
	//file.Write(fltB)
	//file.WriteString("\n")
	//file.Write(fltB)
	//file.Close()
	//file, error = os.Open("/Users/wn/code/6.824-golabs-2020/src/main/tmp")
	//
	//a := test{}
	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	//	line := scanner.Bytes()
	//	json.Unmarshal(line, &a)
	//	fmt.Print(a)
	//	if scanner.Err() != nil {
	//		fmt.Printf(" > Failed with error %v\n", scanner.Err())
	//	}
	//}
	//return

}
