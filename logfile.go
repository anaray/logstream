package logstream

import (
	"bufio"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func parse(file string, reg *regexp.Regexp, control chan struct{}, journal *Journal) int64 {
	//pointer to the file location
	var seekP int64
	var log string
	//var lastLine string
	//check the control channel
	//Open the file
	f, err := os.Open(file)
	if err != nil {
		return -1
	}
	defer func() {
		//create a journal entry and
		//add to the journal
		meta, err := getFileMetaInfo(file)
		if err == nil {
			journalEntry := JournalEntry{
				Signature:   meta.signature,
				File:        file,
				Size:        meta.size,
				ModAt:       meta.modAt,
				Byte_Offset: seekP,
				//Hash:        lastLine,
			}
			journal.Add(journalEntry)
		}
		f.Close()
	}()

	meta, err := getFileMetaInfo(file)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	signature := meta.signature
	je, exists := journal.Get(signature)

	//entry is already there in the journal
	if exists {
		seekP = je.Byte_Offset
	}
	//get
	//seek the correct location in the file
	f.Seek(seekP, 0)
	logscanner := getLogScanner(f, &seekP, reg)
	fmt.Println("Parsing started file: ", file, " at:", time.Now())
L:
	for {
		select {
		case <-control:
			break L
		default:
			more := logscanner.Scan()
			if more {
				log = logscanner.Text()
				fmt.Println(" Parsed Log Entry :=>", log)
			} else {
				break L
			}
		}
	}
	return int64(0)
}

func getLogScanner(file *os.File, seekPos *int64, re *regexp.Regexp) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	logSplitter := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if len(data) > 0 {
			start := re.FindAllIndex(data, 1)
			if len(start) > 0 {
				first := start[0]
				// the pattern is not at start
				if first[0] != 0 {
					*seekPos = *seekPos + int64(first[0])
					return first[0], data[:first[0]], nil
				} else {
					end := re.FindAllIndex(data[first[1]:], 1)
					if len(end) > 0 {
						for _, second := range end {
							if len(second) > 0 {
								adv := second[0] + first[1]
								*seekPos = *seekPos + int64(adv)
								return adv, data[0 : second[0]+first[1]], nil
							}
						}
					} else {
						*seekPos = *seekPos + int64(len(data))
						return len(data), data[0:], nil
					}
				}
			}
		}
		return 0, nil, nil
	}
	scanner.Split(logSplitter)
	return scanner
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
