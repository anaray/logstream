package logstream

import (
	"errors"
)

type LogScrapper struct {
	filepath        string
	delimitor_regex string
}

func NewLogScrapper(logFile string, delimitor string) *LogScrapper {
	ls := LogScrapper{
		filepath:        logFile,
		delimitor_regex: delimitor,
	}
	return &ls
}

func (lgs *LogScrapper) Scrap() error {
	//get the Register and check the journal for this particular file
	if journal == nil {
		return errors.New("logstream: Failed to initialize the scrapper, as is not initialized!")
	}
	return nil
}
