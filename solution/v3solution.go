package solution

import (
	"bufio"
	"os"
)

// V3Solution is an optimized version of V2Solution that replaces the use of bytes.Split
// with a custom byte splitting function.
// execution time: 46s
func V3Solution(filePath string) {
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
		nameB, tempB, _ := CustomByteSplit(row)
		name := string(nameB)
		temp, _ := BytesFloat64ToFloat64V1(tempB)
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
