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
	concurrency  = flag.Int("c", 1, "Number of commands to run in parallel.")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	runtimes := runProgramm(args, *flagCount, *supressOut, *concurrency)
	sort.Slice(runtimes, func(i, j int) bool {
		return runtimes[i] < runtimes[j]

	})

	best := GetBest(runtimes)
	worst := GetWorst(runtimes)
	avg := getAverage(runtimes)

	switch *outputFormat {
	case "m":
		fmt.Printf("Best: %f minutes\n", best.Minutes())
		fmt.Printf("Worst: %f minutes\n", worst.Minutes())
		fmt.Printf("Average: %f minutes\n", avg.Minutes())
	case "s":
		fmt.Printf("Best: %f seconds\n", best.Seconds())
		fmt.Printf("Worst: %f seconds\n", worst.Seconds())
		fmt.Printf("Average: %f seconds\n", avg.Seconds())
	case "ms":
		fmt.Printf("Best: %f milliseconds\n", math.Pow10(3)*best.Seconds())
		fmt.Printf("Worst: %f milliseconds\n", math.Pow10(3)*worst.Seconds())
		fmt.Printf("Average: %f milliseconds\n", math.Pow10(3)*avg.Seconds())
	case "ns":
		fmt.Printf("Best: %f nanoseconds\n", math.Pow10(9)*best.Seconds())
		fmt.Printf("Worst: %f nanoseconds\n", math.Pow10(9)*worst.Seconds())
		fmt.Printf("Average: %f nanoseconds\n", math.Pow10(9)*avg.Seconds())
	default:
		log.Fatal("Unknown output format: ", *outputFormat)
		flag.Usage()
	}
}

func runProgramm(args []string, numRuns int, quiet bool, concurrent int) []time.Duration {
	runtimes := make([]time.Duration, numRuns)
	program := args[0]
	programArgs := args[1:]
	semaphore := make(chan bool, concurrent)

	for i := 0; i < numRuns; i++ {
		// Wait for a free slot
		semaphore <- true
		go func(i int) {
			// Release the slot when we are done
			defer func() { <-semaphore }()
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
		}(i)
	}

	// Wait for all commands to finish
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}

	return runtimes
}

func getAverage(runtimes []time.Duration) time.Duration {
	avg := time.Duration(0)
	for _, r := range runtimes {
		avg += r
	}
	avg /= time.Duration(len(runtimes))
	return avg
}

func GetBest(runtimes []time.Duration) time.Duration {
	best := runtimes[0]
	for _, r := range runtimes {
		if r < best {
			best = r
		}
	}
	return best
}

func GetWorst(runtimes []time.Duration) time.Duration {
	worst := runtimes[0]
	for _, r := range runtimes {
		if r > worst {
			worst = r
		}
	}
	return worst
}
