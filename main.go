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

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

var (
	flagCount    = flag.Int("n", 10, "Number of times to run the command.")
	outputFormat = flag.String("o", "s", "Output format: m (minutes), s (seconds), ms (milliseconds), ns (nanoseconds)")
	supressOut   = flag.Bool("q", false, "Supress output of the command.")
	concurrency  = flag.Int("c", 1, "Number of commands to run in parallel.")
	percentile   = flag.Bool("p", false, "Print percentile values.")
	waitTime     = flag.Duration("w", 0, "Wait time between runs.")
	progress     = flag.Bool("progress", false, "Show progress bar.")
	noColor      = flag.Bool("no-color", false, "Disable colored output.")
	json         = flag.Bool("json", false, "Output results in JSON format.")
	csv          = flag.Bool("csv", false, "Output results in CSV format.")
)

type Result struct {
	Label string
	Value time.Duration
}

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		printUsage()
		return
	}

	// Check if we should use colors
	useColors := !*noColor && isTerminal()

	// Run the program and collect runtimes
	runtimes := runProgramm(args, *flagCount, *supressOut, *concurrency, *waitTime, *progress, useColors)

	// Collect results
	results := []Result{
		{"Best", getBest(runtimes)},
		{"Worst", getWorst(runtimes)},
		{"Average", getAverage(runtimes)},
	}

	if *percentile {
		results = append(results, 
			Result{"Median", getPercentile(runtimes, 50)},
			Result{"90th percentile", getPercentile(runtimes, 90)},
			Result{"95th percentile", getPercentile(runtimes, 95)},
			Result{"99th percentile", getPercentile(runtimes, 99)},
		)
	}

	// Output results in the requested format
	formatString := getFormatString(*outputFormat)
	
	if *json {
		printJSON(results, *outputFormat)
	} else if *csv {
		printCSV(results, *outputFormat)
	} else {
		printResults(results, *outputFormat, formatString, useColors)
	}
}

func printUsage() {
	fmt.Println(colorBold + "best-of" + colorReset + " - Measure execution time of commands")
	fmt.Println("\nUsage:")
	fmt.Println("  best-of [options] -- command [args...]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  best-of -n 5 -- grep -r \"foo\" .")
	fmt.Println("  best-of -o ms -q -- curl https://example.com")
	fmt.Println("  best-of -p -c 10 --progress -- find . -name \"*.go\"")
}

func runProgramm(args []string, numRuns int, quiet bool, concurrent int, waitTime time.Duration, showProgress bool, useColors bool) []time.Duration {
	runtimes := make([]time.Duration, numRuns)
	program := args[0]
	programArgs := args[1:]
	semaphore := make(chan bool, concurrent)

	if showProgress {
		drawProgressBar(0, numRuns, useColors)
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
					if useColors {
						fmt.Printf("%sCommand %s exited with error: %s%s\n", colorRed, program, err, colorReset)
					} else {
						fmt.Printf("Command %s exited with error: %s\n", program, err)
					}
				} else {
					if useColors {
						fmt.Printf("%sFailed to run command: %s%s\n", colorRed, err, colorReset)
					} else {
						fmt.Printf("Failed to run command: %s\n", err)
					}
				}
			}

			runtimes[i] = time.Since(start)
			if showProgress {
				drawProgressBar(i+1, numRuns, useColors)
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

func printResults(results []Result, format string, formatString string, useColors bool) {
	for _, result := range results {
		value := convertTime(result.Value, format)
		
		var labelColor, valueColor string
		if useColors {
			switch result.Label {
			case "Best":
				labelColor = colorGreen + colorBold
				valueColor = colorGreen
			case "Worst":
				labelColor = colorRed + colorBold
				valueColor = colorRed
			case "Average":
				labelColor = colorYellow + colorBold
				valueColor = colorYellow
			case "Median":
				labelColor = colorCyan + colorBold
				valueColor = colorCyan
			default:
				labelColor = colorBlue + colorBold
				valueColor = colorBlue
			}
			
			fmt.Printf("%s%s:%s %s%.6f%s %s\n", 
				labelColor, result.Label, colorReset,
				valueColor, value, colorReset,
				formatString)
		} else {
			fmt.Printf("%s: %.6f %s\n", result.Label, value, formatString)
		}
	}
}

func printJSON(results []Result, format string) {
	fmt.Println("{")
	for i, result := range results {
		value := convertTime(result.Value, format)
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Printf("  \"%s\": %.6f%s\n", result.Label, value, comma)
	}
	fmt.Println("}")
}

func printCSV(results []Result, format string) {
	fmt.Println("Label,Value")
	for _, result := range results {
		value := convertTime(result.Value, format)
		fmt.Printf("%s,%.6f\n", result.Label, value)
	}
}

func drawProgressBar(current int, total int, useColors bool) {
	progressRatio := float64(current) / float64(total)
	progress := int(progressRatio * 100)

	var bar string
	barWidth := 50
	
	// for every 2% we add a character
	for i := 0; i < barWidth; i++ {
		if i < progress/2 {
			if useColors {
				bar += colorGreen + "█" + colorReset
			} else {
				bar += "="
			}
		} else {
			bar += " "
		}
	}

	if useColors {
		fmt.Printf("\r[%s] %s%d%%%s", bar, colorYellow, progress, colorReset)
	} else {
		fmt.Printf("\r[%s] %d%%", bar, progress)
	}
	
	if current == total {
		fmt.Println()
	}
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
