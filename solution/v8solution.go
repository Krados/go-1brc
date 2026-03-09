package solution

import (
	"bytes"
	"context"
	"io"
	"os"
	"runtime"
	"sync"
)

// V8Solution is similar to V7Solution but use fnv1a hash of station name
// instead of station name as key in the map
// execution time: ~2s
func V8Solution(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jobs := make(chan []byte, 1)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	resultChan := make(chan map[uint64]*StationDataV2, 1)
	resultMemo := make(map[string]*StationDataV2)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerFuncV8(ctx, jobs, resultChan)
		}()
	}

	done := make(chan struct{})
	go func() {
		for {
			memo, ok := <-resultChan
			if !ok {
				close(done)
				return
			}
			for _, data := range memo {
				if stationData, exists := resultMemo[data.Name]; exists {
					stationData.Sum += data.Sum
					stationData.Count += data.Count
					if data.Min < stationData.Min {
						stationData.Min = data.Min
					}
					if data.Max > stationData.Max {
						stationData.Max = data.Max
					}
				} else {
					resultMemo[data.Name] = data
				}
			}
		}
	}()

	bufSize := FILE_BUFFER_SIZE
	buffer := make([]byte, bufSize)
	carry := 0
	for {
		n, err := file.Read(buffer[carry:])
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 && err == io.EOF {
			break
		}

		total := carry + n
		cut := bytes.LastIndexByte(buffer[:total], '\n')

		if cut >= 0 {
			end := cut + 1
			data := make([]byte, end)
			copy(data, buffer[:end])
			jobs <- data

			// no need leftover buffer we can just reuse the original buffer
			carry = total - end
			if carry > 0 {
				copy(buffer[:carry], buffer[end:total])
			}
		} else {
			carry = total
			if carry == bufSize {
				panic("single line too long for FILE_BUFFER_SIZE")
			}
		}
	}

	cancel()
	wg.Wait()
	close(jobs)
	close(resultChan)
	<-done
	PrintResultV2(resultMemo)
}
