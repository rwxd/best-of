# best-of

`best-of` lets you check the runtime of program executions.

## Installation

Go to the [releases page](https://github.com/rwxd/best-of/releases) and dowload the latest binary for your system.

Or use `go install`:

```bash
$ go install github.com/rwxd/best-of@latest
```

## Usage

Set the number of executions with the `-n` flag (default 10).

```bash
$ best-of -n 3 -- grep -r "foo" .
Best: 0.031332 seconds
Worst: 0.031558 seconds
Average: 0.031477 seconds
```

Change the output format with the `-o` flag.

```bash
$ best-of -o ms -- grep -r "foo" .
Best: 31.308470 milliseconds
Worst: 31.962246 milliseconds
Average: 31.662080 milliseconds
````

Quiet the output of the programs with the `-q` flag.

```bash
$ best-of -q -- grep -r "foo" .
Best: 0.030725 seconds
Worst: 0.031578 seconds
Average: 0.031138 seconds
```

Use concurrent executions with the `-c` flag (default 1).

```bash
$ best-of -c 10 -- grep -r "foo" .
Best: 0.030725 seconds
Worst: 0.031578 seconds
Average: 0.031138 seconds
```

Get percentiles with the `-p` flag.

```bash
$ best-of -p -q -n 500 -- grep -r "foo" .
Best: 0.030736 s
Worst: 0.032986 s
Average: 0.031279 s
Median: 0.031227 s
90th percentile: 0.031638 s
95th percentile: 0.031779 s
99th percentile: 0.032117 s
```
