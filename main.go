package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/Krados/go-1brc/solution"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var filename = flag.String("file", solution.FILE_NAME, "input file name")
var sv = flag.Int("sv", 7, "solution version to run")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}
	switch *sv {
	case 1:
		solution.V1Solution(*filename)
	case 2:
		solution.V2Solution(*filename)
	case 3:
		solution.V3Solution(*filename)
	case 4:
		solution.V4Solution(*filename)
	case 5:
		solution.V5Solution(*filename)
	case 6:
		solution.V6Solution(*filename)
	case 7:
		solution.V7Solution(*filename)
	default:
		log.Fatalf("unsupported solution version: %d", *sv)
	}
}
