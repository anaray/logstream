package logstream

import (
	"errors"
	"hash/fnv"
	"os"
	"path/filepath"
)

//given a base directory, files are filtered and returned
func getFiles(path, filterPattern string) ([]string, error) {
	fi, err := os.Open(path)
	defer fi.Close()
	if err != nil {
		return nil, err
	}
	finfo, err := fi.Stat()
	if err != nil {
		return nil, err
	}
	var logfiles []string
	if finfo.Mode().IsDir() == false {
		return nil, errors.New("logstream: configured basepath is not a directory!")
	} else {
		logfiles, err = filepath.Glob(path + filterPattern)
		if err != nil {
			return nil, err
		}
	}
	return logfiles, nil
}

// returns 32 bit hash, given a 256 byte content of a file, and
// this acts as a uniue identifier
func getFileSignature(content []byte) uint32 {
	h := fnv.New32a()
	h.Write(content)
	return h.Sum32()
}
