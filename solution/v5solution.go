package solution

import (
	"context"
	"io"
	"os"
	"runtime"
	"sync"
)

// V5Solution is the concurrent version of the v4 solution that uses multiple goroutines to process the file in parallel.
// It reads the file in chunks, sends the chunks to worker goroutines for processing, and then aggregates the results.
// execution time: ~3.6s
func V5Solution(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jobs := make(chan []byte, 1)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	resultChan := make(chan map[string]*StationDataV2, 1)
	resultMemo := make(map[string]*StationDataV2)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerFuncV5(ctx, jobs, resultChan)
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
			for name, data := range memo {
				if stationData, exists := resultMemo[name]; exists {
					stationData.Sum += data.Sum
					stationData.Count += data.Count
					if data.Min < stationData.Min {
						stationData.Min = data.Min
					}
					if data.Max > stationData.Max {
						stationData.Max = data.Max
					}
				} else {
					resultMemo[name] = data
				}
			}
		}
	}()

	bufSize := FILE_BUFFER_SIZE
	buffer := make([]byte, bufSize)
	offset := bufSize
	var hasCut bool
	for {
		hasCut = false
		_, err := file.Read(buffer[bufSize-offset:])
		if err == io.EOF {
			break
		}

		for i := len(buffer) - 1; i >= 0; i-- {
			// 10 = \n
			if buffer[i] == 10 {
				offset = i + 1
				hasCut = true
				break
			}
		}

		if hasCut {
			leftover := make([]byte, bufSize)
			data := make([]byte, offset)
			copy(leftover, buffer[offset:])
			copy(data, buffer[0:offset])
			jobs <- data
			copy(buffer, leftover)
		}

		if err != nil {
			panic(err)
		}
	}

	cancel()
	wg.Wait()
	close(jobs)
	close(resultChan)
	<-done
	PrintResultV2(resultMemo)
}
