package solution

import (
	"bufio"
	"bytes"
	"os"
)

// V2Solution is an optimized version of V1Solution that uses a larger buffer size for the scanner.
// This allows it to read larger chunks of the file at once, reducing the number of I/O operations and improving performance.
// The logic for processing each line remains the same, but the increased buffer size can significantly reduce execution time, especially for large files.
// execution time: 67s
func V2Solution(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	memo := make(map[string]*StationDataV1)
	for scanner.Scan() {
		row := scanner.Bytes()
		data := bytes.Split(row, []byte(";"))
		name := string(data[0])
		temp, _ := BytesFloat64ToFloat64V1(data[1])
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
			memo[name] = &StationDataV1{
				Name:  name,
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		}
	}
	PrintResultV1(memo)
}
