package solution

import (
	"bufio"
	"bytes"
	"os"
)

// V1Solution reads the file line by line, splits the data,
// and updates the station data in a map. It calculates the minimum, maximum,
// and average temperatures for each station and prints the results in a sorted order.
// execution time: 75s
func V1Solution(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
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
