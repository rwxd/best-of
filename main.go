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
	percentile   = flag.Bool("p", false, "Print percentile values.")
	waitTime     = flag.Duration("w", 0, "Wait time between runs.")
	progress     = flag.Bool("progress", false, "Show progress bar.")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	runtimes := runProgramm(args, *flagCount, *supressOut, *concurrency, *waitTime, *progress)

	formatString := getFormatString(*outputFormat)
	printTime("Best", getBest(runtimes), *outputFormat, formatString)
	printTime("Worst", getWorst(runtimes), *outputFormat, formatString)
	printTime("Average", getAverage(runtimes), *outputFormat, formatString)
	if *percentile {
		printTime("Median", getPercentile(runtimes, 50), *outputFormat, formatString)
		printTime("90th percentile", getPercentile(runtimes, 90), *outputFormat, formatString)
		printTime("95th percentile", getPercentile(runtimes, 95), *outputFormat, formatString)
		printTime("99th percentile", getPercentile(runtimes, 99), *outputFormat, formatString)
	}
}

func runProgramm(args []string, numRuns int, quiet bool, concurrent int, waitTime time.Duration, progressBar bool) []time.Duration {
	runtimes := make([]time.Duration, numRuns)
	program := args[0]
	programArgs := args[1:]
	semaphore := make(chan bool, concurrent)

	if progressBar {
		drawProgressBar(0, numRuns)
	}

	for i := 0; i < numRuns; i++ {
		// Wait for a free slot
		semaphore <- true
		go func(i int) {
			// Release the slot when we are done
			defer func() { <-semaphore }()
			time.Sleep(waitTime)
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
				if _, ok := err.(*exec.ExitError); ok {
					fmt.Printf("Command %s exited with error: %s", program, err)
				} else {
					fmt.Printf("Failed to run command: %s", err)
				}
			}

			runtimes[i] = time.Since(start)
			if progressBar {
				drawProgressBar(i+1, numRuns)
			}
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

func getBest(runtimes []time.Duration) time.Duration {
	best := runtimes[0]
	for _, r := range runtimes {
		if r < best {
			best = r
		}
	}
	return best
}

func getWorst(runtimes []time.Duration) time.Duration {
	worst := runtimes[0]
	for _, r := range runtimes {
		if r > worst {
			worst = r
		}
	}
	return worst
}

func getPercentile(runtimes []time.Duration, percentile float64) time.Duration {
	sort.Slice(runtimes, func(i, j int) bool {
		return runtimes[i] < runtimes[j]
	})
	index := (percentile / 100) * float64(len(runtimes))
	if index == float64(int64(index)) {
		return runtimes[int(index)]
	}
	// Interpolation between two points
	lower := runtimes[int(index)]
	if lower == runtimes[len(runtimes)-1] {
		return lower
	}
	upper := runtimes[int(index)+1]
	return lower + time.Duration(float64((upper-lower).Nanoseconds())*(index-float64(int(index))))
}

func convertTime(duration time.Duration, format string) float64 {
	switch format {
	case "m":
		return duration.Minutes()
	case "s":
		return duration.Seconds()
	case "ms":
		return math.Pow10(3) * duration.Seconds()
	case "ns":
		return math.Pow10(9) * duration.Seconds()
	default:
		log.Fatal("Unknown output format: ", format)
		flag.Usage()
		return 0
	}
}

func getFormatString(format string) string {
	switch format {
	case "m":
		return "minutes"
	case "s":
		return "seconds"
	case "ms":
		return "milliseconds"
	case "ns":
		return "nanoseconds"
	default:
		log.Fatal("Unknown output format: ", format)
		flag.Usage()
		return ""
	}
}

func printTime(label string, duration time.Duration, format string, formatString string) {
	fmt.Printf("%s: %f %s\n", label, convertTime(duration, format), formatString)
}

func drawProgressBar(current int, total int) {
	progressRatio := float64(current) / float64(total)
	progress := int(progressRatio * 100)

	var bar string
	// for every 2% we add a "="
	for i := 0; i < 50; i++ {
		if i < progress/2 {
			bar += "="
		} else {
			bar += " "
		}
	}

	fmt.Printf("\r[%s] %d%%", bar, progress)
	if current == total {
		fmt.Println()
	}
}
