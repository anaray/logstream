package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"logstream"
	"os"
	"time"
)

type Conf struct {
	Path          string `json:"base_directory"`
	FilterPattern string `json:"filter_pattern"`
	LogDelimRegex string `json:"log_delim_regex"`
	Interval      int64  `json:"interval"`
	Timeout       int64  `json:"timeout"`
	JournalPath   string `json:"journal_path"`
}

func main() {
	logger := logstream.Logger(os.Stdout)
	logger.Logf("Initializing Logstream: 0.1")
	f := flag.String("conf", ".", "the configuration file")
	flag.Parse()
	c, err := ioutil.ReadFile(*f)
	if err != nil {
		panic(fmt.Sprintf("logstream: failed to read configuration file %s\n", *f))
		os.Exit(1)
	}
	conf := Conf{}
	json.Unmarshal(c, &conf)
	logger.Logf("Creating LogStream Agent ...")
	logger.Logf("File Path: %s\n", conf.Path)
	logger.Logf("Filter: %s\n", conf.FilterPattern)
	logger.Logf("Log Delim Regex: %s\n", conf.LogDelimRegex)
	logger.Logf("Gather Interval: %d\n", conf.Interval)
	logger.Logf("Gather Timeout: %d\n", conf.Timeout)
	logger.Logf("Journal Path: %s\n", conf.JournalPath)
	agent := logstream.NewAgent(conf.Path,
		conf.FilterPattern,
		conf.LogDelimRegex,
		time.Duration(conf.Interval)*time.Second,
		time.Duration(conf.Timeout)*time.Second,
		conf.JournalPath)
	//create
	agent.Start()
}
