package main

import (
	"os"
	"path"
)

// ExportFile contains the information required for Exporter to manipulate files during export.
type ExportFile struct {
	descriptor *os.File
	name       string
	format     FormatInterface
}

// ExportInterface provides a mock target for Export objects
type ExportInterface interface {
	Format() FormatInterface
	GetFileName() string
	GetSource() string
}

// Export is a parent object which contains the source of the export and the object to be exported
type Export struct {
	source  string
	destDir string
	obj     FormatInterface
}

// GetFileName returns the path to which the object should be exported
func (t *Export) GetFileName() string {
	baseFilename := path.Base(t.source)
	var extension = path.Ext(baseFilename)
	var fileName = baseFilename[0 : len(baseFilename)-len(extension)]
	return path.Join(t.destDir, fileName+t.obj.GetFileExtension())
}

// Format returns the FormatInterface object which provides information of object serialization
func (t *Export) Format() FormatInterface {
	return t.obj
}

// GetSource returns the the source file of the export
func (t *Export) GetSource() string {
	return t.source
}

// Exporter streams objects from memory to disk, record by record, terminating and
// closing the file once all records have been processed.
func Exporter(
	assignedFile *string, // Empty string is idle
	exportStream chan ExportInterface,
	processingContext ProcessingContextInterface,
	killChan chan bool,
	errChan chan error,
) {
	openExports := make([]ExportFile, 0)
	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			cleanUpExportFiles(openExports, errChan)
		case export := <-exportStream:
			// fmt.Println("Exporter Recieved", export.GetFileName())
			handled := false
			fileName := export.GetFileName()
			for _, val := range openExports {
				if val.name == fileName {
					// fmt.Println("Exporter Wrote", export.GetFileName())
					if _, err := val.descriptor.WriteString(
						export.Format().StandardRecord(),
					); err != nil {
						errChan <- err
					}
					handled = true
					break
				}
			}

			if !handled {
				newFile, cErr := os.Create(fileName)
				if cErr != nil {
					errChan <- cErr
					continue
				}
				// fmt.Println("Exporter Opened", export.GetFileName())
				_, wErr := newFile.WriteString(export.Format().FirstRecord())
				if wErr != nil {
					newFile.Close()
					errChan <- wErr
					continue
				}
				// fmt.Println("Exporter Wrote", export.GetFileName())
				openExports = append(openExports, ExportFile{newFile, fileName, export.Format()})
			}

			if err := processingContext.IncrementRecordsProcessed(*assignedFile); err != nil {
				errChan <- err
			}

			complete, pcErr := processingContext.ProcessingComplete(*assignedFile)
			if pcErr != nil {
				errChan <- pcErr
			}

			if complete {
				cleanUpExportFiles(openExports, errChan)

				openExports = make([]ExportFile, 0)
				*assignedFile = ""
			}
		}
	}
}

func cleanUpExportFiles(openExports []ExportFile, errChan chan error) {
	for _, val := range openExports {
		_, wErr := val.descriptor.WriteString(val.format.TerminatingString())
		if wErr != nil {
			errChan <- wErr
		}
		val.descriptor.Close()
		// fmt.Println("Exporter Closed", val.name)
	}
}
