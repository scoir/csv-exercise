package main

import "errors"

// DirMonitor keeps track of what files are where and is responsible for determining which files are
// processed. It also flags files for deletion.
func DirMonitor(inDir, outDir, errDir *string, toProcess, toDelete chan string) error {
	return errors.New("unimplemented")
}
