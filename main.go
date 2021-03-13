package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	. "rabbits_wolfs/animals"
	. "rabbits_wolfs/field"
	"runtime"
	"runtime/pprof"
	"time"
)

const Epoch = 100

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	fmt.Print("Time to live (", Epoch, "):")

	data4Loop := new(loopData)

	fmt.Scanln(&data4Loop.time2Live)
	if data4Loop.time2Live == 0 {
		data4Loop.time2Live = Epoch
	}

	fmt.Print("field size (", Size, "):")
	var fieldSize int
	fmt.Scanln(&fieldSize)
	if fieldSize == 0 {
		fieldSize = Size
	}
	fmt.Println("field size'll be ", fieldSize, "x", fieldSize, "(", fieldSize*fieldSize, ")")

	data4Loop.field = NewField(fieldSize)

	fmt.Print("Rabbits start count (", RabbitDefCount, "):")
	var rabbitsCount int
	fmt.Scanln(&rabbitsCount)
	if rabbitsCount == 0 {
		rabbitsCount = RabbitDefCount
	}

	fmt.Print("Wolfs start count (", WolfDefCount, "):")
	var wolfsCount int
	fmt.Scanln(&wolfsCount)
	if wolfsCount == 0 {
		wolfsCount = WolfDefCount
	}

	rand.Seed(time.Now().UnixNano())
	data4Loop.animals = []Animals{NewRabbits(rabbitsCount), NewWolfs(wolfsCount)}
	start := time.Now()
	loop(data4Loop)
	fmt.Println(time.Since(start))

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}