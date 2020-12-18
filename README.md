# SCOIR Technical Interview for Back-End Engineers
This repo contains an exercise intended for Back-End Engineers.

## Instructions
1. Fork this repo.
1. Using technology of your choice, complete [the assignment](./Assignment.md).
1. Update this README with
    * a `How-To` section containing any instructions needed to execute your program.
    * an `Assumptions` section containing documentation on any assumptions made while interpreting the requirements.
1. Before the deadline, submit a pull request with your solution.

## Expectations
1. Please take no more than 8 hours to work on this exercise. Complete as much as possible and then submit your solution.
1. This exercise is meant to showcase how you work. With consideration to the time limit, do your best to treat it like a production system.

## Assumptions
* Files which have been converted and removed from the output directory will not be re-added to the input directory unless the user wants them re-converted.
    - The program is expected to be naive -- every file in input upon start and added to input during lifetime will be processed.
* The “name collision” rule for CSV validation error files applies universally.
    - e.g. if there are files with identical names in output and input, overwrite the output with the file in input.
* If there is an OS-level file-read interruption of the daemon, it should log the error and continue. This could result in a “partial” conversion of a file, where the final error in the error log for conversion would be “FILE,<err_msg>”. After this, the file will be considered “processed” and deleted.
    - e.g. A file is moved, modified, deleted, corrupted, etc
    - NOTE: If the file is flagged for “reprocessing” due to reception of a fileModified event, then it would be reprocessed.
* If there is a parsing error encountered when reading a CSV file, treat it as another verification error and log it in the same manner.
* If there is an OS-level file-write interruption of the daemon, it should log the error and shutdown immediately, without deleting anything in input that is still being processed.
    - If I can’t write out the conversion, there is no point in continuing, I’m just deleting the input files without writing anything out.
* The Go runtime can handle load distribution on different channels given the processes to do so.
    - A “naive” load distribution approach because something more advanced would take up too much time.
* The filesystem targeted will not be FUSE, which does not work with the fsnotify library utilized to track directory changes.
* If the program experiences an internal error processing a file, that error will be logged and the file will not be deleted.
* The built-in CSV and JSON packages in GoLang will handle any characters which need to be escaped or otherwise modified.
* If the header for the input csv file is incorrect, stop processing the file.

## How-To
You can run it using `go run .` or build an executable using `go build .` and run that. I'll include the usage message below, although I didn't get the chance to hook all of them up. Right now I'd stick to the basic required args and assume nothing else is hooked up.
```
  -c    Runs the program in 'careful' mode, requesting additional user input
  -d    Runs the program in multi-threaded mode
  -e string
        The error directory
  -h    Displays a usage message
  -i string
        The input directory
  -l string
        The file logs are written to
  -mf int
        Controls the size of the file processing queue the program to maintain. Must equal or exceed the maximum amount of files being processed at once. (default 100)
  -mr int
        Controls the amount of queued records each channel is allowed to main. Should correspond with the maximum amount of records per file. (default 1000)
  -o string
        The output directory
  -v    Increases reporting when present
```

## Current Program State
The basic pipeline has been assembled, but testing is not complete and there are definitely bugs -- I found one to do with tracking file record completion. More unit testing is needed, in addition to automated full-stack tests. Reporting and Logging also need a lot of work. Current implementation does not adhere to the assumptions around error-handling, with error channels having no listener.

There were a few design iterations which left some inconsistencies between services, resulting in some main-level converters, spliters, and utility providers.

## Program Structure
The program is composed of several "scale-at-will" services along with a few singleton services to coordinate everything. As of now, there is a hard limit on concurrently processesing files linked to the number of available processing units. During design I chose between this and finding a way to mititage kernel-level call churn, which I determined would take too long for an 8 hour project. Right now the Go Runtime is in charge of juggling which goroutines get processor time, with the amount of service instances limited to a maximum of the number of available CPUs.
* DirMonitor is a singleton and reports any activity in the input directory, as well as scanning it upon initialization
* CSVImporter is scale-at-will, moving files from disk to memory record-by-record, with one process per file.
* Verifier is scale-at-will, taking in records and verifying them.
* ExportManager is a singleton, and manages the exporting instances. It relies on ProcessingContext to ensure there are never so many files in the pipeline that it gets overwhelmed. It streams objects from memory to disk record-by-record, with one input file per Exporter.