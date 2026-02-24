package solution

import (
	"bufio"
	"os"
)

// V4Solution is the optimized version of the v3 solution that replaces the use of strconv.ParseFloat
// with a custom function that converts the byte slice directly to an int64 representation
// of the float value multiplied by 10. This avoids the overhead of parsing the float and allows us
// to work with integers, which is more efficient for our use case.
// execution time: 36s
func V4Solution(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	memo := make(map[string]*StationDataV2)
	for scanner.Scan() {
		row := scanner.Bytes()
		nameB, tempB, _ := CustomByteSplit(row)
		name := string(nameB)
		temp := ByteFloat64ToInt64V1(tempB)
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
	PrintResultV2(memo)
}
