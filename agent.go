package logstream

import (
	"fmt"
	"time"
)

type Agent struct {
	gatherInterval time.Duration
	gatherTimeout  time.Duration
	shutdown       chan struct{}
	basePath       string
	filterPattern  string
	journalPath    string
}

func NewAgent(basePath string, filterPattern string, interval, timeout time.Duration, journalPath string) *Agent {
	agent := Agent{
		gatherInterval: time.Duration(interval),
		gatherTimeout:  time.Duration(timeout),
		shutdown:       make(chan struct{}),
		basePath:       basePath,
		filterPattern:  filterPattern,
		journalPath:    journalPath,
	}
	return &agent
}

func (agent *Agent) Start() {
	//initialize a journal, with gob to serialize to disk
	journal, journal_chan, sweep_chan, err := CreateJournal(agent.journalPath, writeToGob, loadFromGob)
	if err != nil {
		panic(err)
	}
	fmt.Println(journal, journal_chan, sweep_chan, err)

	ticker := time.NewTicker(agent.gatherInterval)
	defer ticker.Stop()
	for {
		files, _ := getFiles(agent.basePath, agent.filterPattern)
		//timeout is a optimistic way to divide the parsing among the
		//files. It is equally divided among the files
		timeOutInSeconds := agent.gatherTimeout.Seconds() / float64(len(files))
		for _, file := range files {
			meta, err := getFileMetaInfo(file)
			if err != nil {
				fmt.Println(err)
				continue
			}
			f := meta.signature
			fmt.Println("File Signature :", f)
			entry, ok := journal.Entries[f]
			fmt.Println("Entry :", entry, ok)
			gather(file, time.Duration(timeOutInSeconds)*time.Second, agent.shutdown)
			// write to details to journal for this entry in a goroutine
		}
		select {
		case <-agent.shutdown:
			return
		case <-ticker.C:
			fmt.Println("Calling Gather .....", time.Now())
			continue
		}
	}
}

func gather(file string, timeout time.Duration, shutdown chan struct{}) {
	timeOutTicker := time.NewTicker(timeout)
	defer timeOutTicker.Stop()
	done := make(chan error)
	control := make(chan struct{})
	go func() {
		done <- parse(file, control)
	}()
	for {
		select {
		case <-done:
			fmt.Println("parsing completed !!!!", time.Now())
			return
		case <-shutdown:
			fmt.Println("shutdown called !")
			return
		case <-timeOutTicker.C:
			fmt.Println("timeout called! stopping the parsing of %s", file)
			//close(control)
			var stop struct{}
			control <- stop
			//continue
			return
		}
	}
}
