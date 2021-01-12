package parser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/8bitstout/csv-exercise/record"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Parser holds the structure defining a parser object
type Parser struct {
	RowLength       int
	InputDirectory  string
	OutputDirectory string
	ErrorsDirectory string
}

// Watch listens to the input-directory for new files and
// processes the file in a separate Goroutine. If any errors
// are found while parsing a csv file, they will be outputted to
// a csv file in the errors-directory. Csv files will be parsed to
// a JSON file and outputted to the output-directory
func (p *Parser) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event: ", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					fileName := filepath.Base(event.Name)
					fileName = strings.TrimSuffix(fileName, path.Ext(fileName))
					r, errorLog := p.ParseRecordsToJSON(event.Name)

					if errorLog.HasData() {
						fmt.Println(errorLog.Errors)
						p.WriteErrorsToFile(fileName, errorLog.Errors)
					}

					p.WriteRecordsToJSON(fileName, r)
					err := os.Remove(event.Name)

					if err != nil {
						log.Println("Error removing parsed CSV:", err)
					}
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(err)
			}
		}
	}()

	err = watcher.Add(p.InputDirectory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// ParseRecordsToJSON takes a file path, parses the csv to JSON, and returns
// an array of Record objects to be written to a JSON file and an ErrorLog object
// that contains any found errors to be written to a csv file
func (p *Parser) ParseRecordsToJSON(file string) ([]*record.Record, *record.ErrorLog) {
	var records []*record.Record
	errorLog := record.NewErrorLog()
	csvfile, err := os.Open(file)

	if err != nil {
		log.Fatal("Error opening file:", err)
	}

	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	header, err := r.Read()

	if err != nil {
		log.Fatal(err)
	}

	if len(header) < p.RowLength {
		log.Fatalf("Error processing csv file: header row is expected to have %v elements", p.RowLength)
	}

	for i, h := range header {
		if _, ok := record.CsvHeader[h]; !ok {
			log.Fatal(h, "is not a valid header name")
		}

		if i != record.CsvHeader[h] {
			log.Fatalf("Header \"%s\" is in the wrong position. Expected position %v, but found position %v", h, record.CsvHeader[h], i)
		}
	}

	rowNumber := 1

	for {
		var rowHasError bool
		rowNumberStr := strconv.Itoa(rowNumber)

		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			rowHasError = true
			errorLog.Append(rowNumberStr, fmt.Sprint(err))
		}

		if len(row) < p.RowLength {
			rowHasError = true
			errorLog.Append(rowNumberStr, fmt.Sprintf("Error processing csv file: row is expected to have %v elements", p.RowLength))
		}

		id, err := strconv.Atoi(row[0])

		if err != nil {
			rowHasError = true
			errorLog.Append(rowNumberStr, fmt.Sprint("Expected a positive integer for the ID field, but received an unexpected value"))
		}

		firstName := row[1]
		middleName := row[2]
		lastName := row[3]
		phoneNum := row[4]

		n, err := record.NewName(firstName, middleName, lastName)

		if err != nil {
			rowHasError = true
			errorLog.Append(rowNumberStr, fmt.Sprint(err))
		}

		pn, err := record.NewPhoneNumber(phoneNum)

		if err != nil {
			rowHasError = true
			errorLog.Append(rowNumberStr, fmt.Sprint(err))
		}

		r, err := record.NewRecord(uint32(id), n, pn)

		if err != nil {
			rowHasError = true
			errorLog.Append(rowNumberStr, fmt.Sprint(err))
		}

		if !rowHasError {
			records = append(records, r)
		}
		rowNumber++
	}

	return records, errorLog
}

// WriteRecordsToJSON takes a filename and an array of Record object pointers
// and writes the Record objects to a JSON file
func (p *Parser) WriteRecordsToJSON(file string, records []*record.Record) {
	data, _ := json.MarshalIndent(records, "", "")
	err := ioutil.WriteFile(p.OutputDirectory+"/"+file+".json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// WriteErrorsToFile takes a filename and an array of error logs
// and writes the error logs to a csv files
func (p *Parser) WriteErrorsToFile(file string, errors [][]string) {
	f, err := os.Create(fmt.Sprintf("%s/%s.csv", p.ErrorsDirectory, file))
	defer f.Close()
	if err != nil {
		log.Fatal("Could not create file:", err)
	}

	w := csv.NewWriter(f)

	for _, row := range errors {
		_ = w.Write(row)
	}

	w.Flush()
}

// NewParser returns a pointer to a new Parser object
func NewParser(rowLength int, inputDirectory, outputDirectory, errorsDirectory string) *Parser {
	directories := []string{inputDirectory, outputDirectory, errorsDirectory}

	for _, d := range directories {
		_, err := os.Stat(d)
		if os.IsNotExist(err) {
			err = os.MkdirAll(d, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return &Parser{
		RowLength:       rowLength,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		ErrorsDirectory: errorsDirectory,
	}
}
