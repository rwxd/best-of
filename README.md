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
Best: 0.030763 seconds
Worst: 0.034639 seconds
Average: 0.031361 seconds
Median: 0.031313 seconds
90th percentile: 0.031742 seconds
95th percentile: 0.031873 seconds
99th percentile: 0.032652 seconds
```

Wait between runs with the `-w` flag (default 0).

```bash
best-of -w 3s -- grep -r "foo" .
```
