package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"time"
)

var (
	flagCount    = flag.Int("n", 10, "Number of times to run the command.")
	outputFormat = flag.String("o", "s", "Output format: m (minutes), s (seconds), ms (milliseconds), ns (nanoseconds)")
	supressOut   = flag.Bool("q", false, "Supress output of the command.")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	runtimes := runProgramm(args, *flagCount, *supressOut)
	sort.Slice(runtimes, func(i, j int) bool {
		return runtimes[i] < runtimes[j]

	})

	best := runtimes[0]

	switch *outputFormat {
	case "m":
		fmt.Printf("Best time: %f minutes\n", best.Minutes())
	case "s":
		fmt.Printf("Best time: %f seconds\n", best.Seconds())
	case "ms":
		fmt.Printf("Best time: %f milliseconds\n", math.Pow10(3)*best.Seconds())
	case "ns":
		fmt.Printf("Best time: %f nanoseconds\n", math.Pow10(9)*best.Seconds())
	default:
		log.Fatal("Unknown output format: ", *outputFormat)
	}
}

func runProgramm(args []string, numRuns int, quiet bool) []time.Duration {
	runtimes := make([]time.Duration, numRuns)
	program := args[0]
	programArgs := args[1:]

	for i := 0; i < numRuns; i++ {
		start := time.Now()
		cmd := exec.Command(program, programArgs...)

		if quiet {
			cmd.Stdout = io.Discard
			cmd.Stderr = io.Discard
		} else {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			log.Fatal("Failed to run command: ", err)
		}
		runtimes[i] = time.Since(start)
	}
	return runtimes
}
