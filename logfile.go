package logstream

import (
	"errors"
	"hash/fnv"
	"os"
	"path/filepath"
	"time"
)

func parse(file string) error {
  //f, err := os.Open(file)
  //defer f.Close()
  time.Sleep(20 * time.Second)
  return nil
}

type LogFileMetaInfo struct {
	signature uint32
	size      int64
	modAt     time.Time
}

const BUFFER_LEN = 256

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

func getFileMetaInfo(file string) (*LogFileMetaInfo, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	s, err := getFileSignature(f)
	if err != nil {
		return nil, err
	}
	fi, ferr := f.Stat()
	if ferr != nil {
		return nil, ferr
	}
	lgmi := LogFileMetaInfo{
		signature: s,
		size:      fi.Size(),
		modAt:     fi.ModTime(),
	}
	return &lgmi, nil
}

// returns 32 bit hash, this acts as a uniue identifier
func getFileSignature(file *os.File) (uint32, error) {
	b := make([]byte, BUFFER_LEN)
	read, err := file.Read(b)
	if err != nil {
		return 0, err
	}
	// not enough data to create a signature,
	// will try to create in next iteration
	if read < BUFFER_LEN {
		return 0, errors.New("logstream: file < 256 bytes, next time!")
	}
	h := fnv.New32a()
	h.Write(b)
	return h.Sum32(), nil
}
