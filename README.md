# SCOIR Technical Interview
This repo contains an exercise intended for engineers.

## Instructions
1. Fork this repo.
1. Using technology of your choice, complete [the assignment](./Assignment.md).
1. Update this README with
    * a `How-To` section containing any instructions needed to execute your program.
    * an `Assumptions` section containing documentation on any assumptions made while interpreting the requirements.
1. Before the deadline, submit a pull request with your solution.

## Expectations
1. Please take no more than 8 hours to work on this exercise. Complete as much as possible and then submit your solution.
2. This exercise is meant to showcase how you work. With consideration to the time limit, do your best to treat it like a production system.

## How-To
### Requirements
- Go version 1.15.6
### Dependencies
- [fsnotify](https://github.com/fsnotify/fsnotify) for watching for changes in a given directory
- [cli](https://github.com/urfave/cli) for providing a cli interface
## Running
This program can be run with either of the following methods
```
go run ./cmd/csv-exercise/csv-exercise.go
```
```
go build ./cmd/csv-exercise/csv-exercise.go
csv-exercise
```

By default, this program will create an input, output, and errors directory.
You can define your own directories by using any of the following flags
```
csv-exercise [--input-dir ${input_directory}] [--output-dir ${output_directory}] [--errors-dir ${errors_directory}]
```

There is a shell script to create mock data. This script should be called after the Go program is running.
The shell script has the input directory hardcoded to match the default input directory in the program.
If you call the program with flags to provide your own input directory, the shell script must be updated
to use that directory.
```
# Omit the extension from ${file_name}
./generate_csv.sh ${file_name}
```
## Assumptions
- Mock data was hardcoded and can be generated with the provided shell script.
- This program has been tested on MacOS.
- The program will terminate if a csv with invalid headers is provided as input.
I believe this is a reasonable behavior given that an invalid format would affect
every record in the file - we should not log this and instead terminate.
- The program will also terminate if it is unable to produce a required file (error log or json)
- The program will not attempt to clean inputs of each field such as whitespaces or numbers
where characters are expected to be such as in a name.
- The Parser object contains a property that is a map for keeping track of processed and unprocessed files.
The assumption here is that the map is accessed in only one Goroutine so we will not worry about write access.
Otherwise we need to consider thread saftey.
- Even if a file is set as processed in the processedFiles map, we will still process it again and overwrite
the previous file given the requirement to overwrite files when name collisions occur.
- No actions are taken if any directories are removed during program runtime