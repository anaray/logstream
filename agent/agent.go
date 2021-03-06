package logstream

import (
	"fmt"
	"regexp"
	"time"
)

type Agent struct {
	gatherInterval time.Duration
	gatherTimeout  time.Duration
	shutdown       chan struct{}
	Output         chan Result
	basePath       string
	filterPattern  string
	regexDelim     string
	logType        string
	journalPath    string
}

type Result struct {
	File     string            `json:"file"`
	Log      string            `json:"file"`
	Type     string            `json:"file"`
	ParsedAt time.Time         `json:"file"`
	MetaInfo map[string]string `json:"meta_info"`
}

const OUTPUT_LEN = 256

func NewAgent(basePath, filterPattern, regexDelim, logType string, interval, timeout time.Duration, journalPath string) *Agent {
	agent := Agent{
		gatherInterval: time.Duration(interval),
		gatherTimeout:  time.Duration(timeout),
		shutdown:       make(chan struct{}),
		Output:         make(chan Result, OUTPUT_LEN),
		basePath:       basePath,
		filterPattern:  filterPattern,
		regexDelim:     regexDelim,
		logType:        logType,
		journalPath:    journalPath,
	}
	return &agent
}

func (agent *Agent) Start() {
	//initialize a journal, with gob to serialize to disk
	//journal, journal_chan, sweep_chan, err := CreateJournal(agent.journalPath, writeToGob, loadFromGob)
	//keep the log file delim marker regex pattern compiled
	delim_regex := regexp.MustCompile(agent.regexDelim)
	ticker := time.NewTicker(agent.gatherInterval)
	defer ticker.Stop()
	for {
		journal, err := GetJournal(agent.journalPath, loadFromGob)
		if err != nil {
			panic(err)
		}
		files, _ := getFiles(agent.basePath, agent.filterPattern)
		fmt.Println("Selected Files :", files)
		//timeout is a optimistic way to divide the parsing among the
		//files. It is equally divided among the files
		timeOutInSeconds := agent.gatherTimeout.Seconds() / float64(len(files))
		for _, file := range files {
			gather(file, delim_regex, agent.Output, agent.logType, time.Duration(timeOutInSeconds)*time.Second, agent.shutdown, journal)
			// write to details to journal for this entry in a goroutine
		}
		journal.Sweep(writeToGob)
		select {
		case <-agent.shutdown:
			return
		case <-ticker.C:
			continue
		}
	}
}

func gather(file string, delim *regexp.Regexp, output chan Result, logType string, timeout time.Duration, shutdown chan struct{}, journal *Journal) {
	//done := make(chan error)
	done := make(chan int64)
	control := make(chan struct{}, 1)
	var t0 time.Time
	timeOutTicker := time.NewTicker(timeout)
	defer timeOutTicker.Stop()

	go func() {
		t0 = time.Now()
		done <- parse(file, delim, logType, output, control, journal)
	}()

	for {
		select {
		case <-done:
			//update the file seek pos
			fmt.Println("parsing completed !!!! file:", file)
			return
			break
		case <-shutdown:
			fmt.Println("shutdown called !")
			break
			return
		case <-timeOutTicker.C:
			fmt.Println("timeout called! stopping the parsing of %s", file)
			type stop struct{}
			s := stop{}
			control <- s
			t1 := time.Now()
			fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
		}
	}
	return
}
