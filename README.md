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
Best time: 0.001537 seconds
```

Change the output format with the `-o` flag.

```bash
$ best-of -o ms -- grep -r "foo" .
Best time: 1.558415 milliseconds
````

Quiet the output of the programs with the `-q` flag.

```bash
$ best-of -q -- grep -r "foo" .
Best time: 0.001537 seconds
```
