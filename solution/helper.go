package solution

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

const FILE_BUFFER_SIZE = 1024 * 1024
const SCANNER_BUFFER_SIZE = 1024 * 1024 * 8
const FILE_NAME_10 = "measurements-10.txt"
const FILE_NAME_10M = "measurements-10000000.txt"
const FILE_NAME = "measurements.txt"

var POW_ARY = []int64{1, 10, 100, 1000}

type StationDataV1 struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

type StationDataV2 struct {
	Name  string
	Min   int64
	Max   int64
	Sum   int64
	Count int
}

func BytesFloat64ToFloat64V1(b []byte) (float64, error) {
	return strconv.ParseFloat(string(b), 64)
}

func CustomByteSplit(b []byte) ([]byte, []byte, bool) {
	for i, v := range b {
		// 59 = ';'
		if v == 59 {
			return b[:i], b[i+1:], true
		}
	}

	return nil, nil, false
}

func CustomByteSplitIdx(b []byte) (int, bool) {
	for i, v := range b {
		// 59 = ';'
		if v == 59 {
			return i, true
		}
	}

	return -1, false
}

func ByteFloat64ToInt64V1(b []byte) int64 {
	isNegative := false
	hasDot := false
	var res int64 = 0
	pow := 0
	for i := len(b) - 1; i >= 0; i-- {
		// 45 = '-'
		if b[i] == 45 {
			isNegative = true
			continue
		}
		// 46 = '.'
		if b[i] == 46 {
			hasDot = true
			continue
		}
		var s int64 = 1
		for i := 0; i < pow; i++ {
			s *= 10
		}
		res += int64(b[i]-48) * s
		pow++
	}
	if !hasDot {
		res *= 10
	}
	if isNegative {
		res = -res
	}
	return res
}

func ByteFloat64ToInt64V2(b []byte) int64 {
	isNegative := false
	hasDot := false
	var res int64 = 0
	pow := 0
	for i := len(b) - 1; i >= 0; i-- {
		// 45 = '-'
		if b[i] == 45 {
			isNegative = true
			continue
		}
		// 46 = '.'
		if b[i] == 46 {
			hasDot = true
			continue
		}
		res += int64(b[i]-48) * POW_ARY[pow]
		pow++
	}
	if !hasDot {
		res *= 10
	}
	if isNegative {
		res = -res
	}
	return res
}

func PrintResultV1(data map[string]*StationDataV1) {
	result := make(map[string]*StationDataV1, len(data))
	keys := make([]string, 0, len(data))
	for _, v := range data {
		keys = append(keys, v.Name)
		result[v.Name] = v
	}
	sort.Strings(keys)

	print("{")
	for _, k := range keys {
		v := result[k]
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", k, v.Min, v.Sum/float64(v.Count), v.Max)
	}
	print("}\n")
}

func PrintResultV2(data map[string]*StationDataV2) {
	result := make(map[string]*StationDataV2, len(data))
	keys := make([]string, 0, len(data))
	for _, v := range data {
		keys = append(keys, v.Name)
		result[v.Name] = v
	}
	sort.Strings(keys)

	print("{")
	for _, k := range keys {
		v := result[k]
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", k, float64(v.Min)/10, float64(v.Sum)/float64(v.Count)/10, float64(v.Max)/10)
	}
	print("}\n")
}

func workerFuncV5(ctx context.Context, jobs chan []byte, resultChan chan map[string]*StationDataV2) {
	buf := make([]byte, SCANNER_BUFFER_SIZE)
	for {
		select {
		case b := <-jobs:
			memo := make(map[string]*StationDataV2)
			scanner := bufio.NewScanner(bytes.NewReader(b))
			scanner.Buffer(buf, SCANNER_BUFFER_SIZE)
			for scanner.Scan() {
				row := scanner.Bytes()
				nameByte, tempBytes, ok := CustomByteSplit(row)
				if !ok {
					log.Fatalln("Error parsing row: expected ';' separator not found")
				}
				name := string(nameByte)
				temp := ByteFloat64ToInt64V1(tempBytes)
				if stationData, exists := memo[name]; exists {
					stationData.Sum += temp
					stationData.Count++
					if temp < stationData.Min {
						stationData.Min = temp
					}
					if temp > stationData.Max {
						stationData.Max = temp
					}
				} else {
					memo[name] = &StationDataV2{
						Name:  name,
						Min:   temp,
						Max:   temp,
						Sum:   temp,
						Count: 1,
					}
				}
			}
			// send the result back to the main goroutine
			resultChan <- memo
		case <-ctx.Done():
			return
		}
	}
}

func workerFuncV7(ctx context.Context, jobs chan []byte, resultChan chan map[string]*StationDataV2) {
	buf := make([]byte, SCANNER_BUFFER_SIZE)
	for {
		select {
		case b := <-jobs:
			memo := make(map[string]*StationDataV2)
			scanner := bufio.NewScanner(bytes.NewReader(b))
			scanner.Buffer(buf, SCANNER_BUFFER_SIZE)
			for scanner.Scan() {
				row := scanner.Bytes()
				nameByte, tempBytes, ok := CustomByteSplit(row)
				if !ok {
					log.Fatalln("Error parsing row: expected ';' separator not found")
				}
				name := string(nameByte)
				temp := ByteFloat64ToInt64V2(tempBytes)
				if stationData, exists := memo[name]; exists {
					stationData.Sum += temp
					stationData.Count++
					if temp < stationData.Min {
						stationData.Min = temp
					}
					if temp > stationData.Max {
						stationData.Max = temp
					}
				} else {
					memo[name] = &StationDataV2{
						Name:  name,
						Min:   temp,
						Max:   temp,
						Sum:   temp,
						Count: 1,
					}
				}
			}
			// send the result back to the main goroutine
			resultChan <- memo
		case <-ctx.Done():
			return
		}
	}
}

// TestReadBufioScanner is a helper function to read the file using bufio.Scanner
// and simply iterate through the lines without doing any processing.
// execution time: ~10s
func TestReadBufioScanner(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_ = scanner.Bytes()
	}
}

// TestReadBufioScannerWithBuffer is a helper function to read the file using bufio.Scanner
// with a custom buffer size
// execution time: ~6s
func TestReadBufioScannerWithBuffer(filePath string, bufferSize int) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, bufferSize), bufferSize)
	for scanner.Scan() {
		_ = scanner.Bytes()
	}
}

// TestReadFileRead is a helper function to read the file using os.File.Read
// with a custom buffer size
// execution time: ~1.2s
func TestReadFileRead(filePath string, bufferSize int) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buffer := make([]byte, bufferSize)
	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
	}
}
