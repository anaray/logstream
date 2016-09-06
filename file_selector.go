package logstream

import (
	"errors"
	"os"
	"path/filepath"
)

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
		return nil, errors.New("logstream: configured basepath is not a directory !")
	} else {
		logfiles, err = filepath.Glob(path + filterPattern)
		if err != nil {
			return nil, err
		}
	}
	return logfiles, nil
}
