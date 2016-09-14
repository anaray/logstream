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
}

func NewAgent(basePath string, filterPattern string, interval, timeout time.Duration) *Agent {
	agent := Agent{
		gatherInterval: time.Duration(interval),
		gatherTimeout:  time.Duration(timeout),
		shutdown:       make(chan struct{}),
		basePath:       basePath,
		filterPattern:  filterPattern,
	}
	return &agent
}

func (agent *Agent) Start() {
	//CreateJournal()
	ticker := time.NewTicker(agent.gatherInterval)
	defer ticker.Stop()
	for {
		files, _ := getFiles(agent.basePath, agent.filterPattern)
		//timeout is a optimistic way to divide the parsing among the
		//files. It is equally divided among the files
		//fmt.Println("1 :", agent.gatherTimeout.Seconds())
		//fmt.Println("2 :", float64(len(files)))
		timeOutInSeconds := agent.gatherTimeout.Seconds() / float64(len(files))
		//fmt.Println("Calculated timeout :", timeOutInSeconds)

		for _, file := range files {
			gather(file, time.Duration(timeOutInSeconds)*time.Second, agent.shutdown)
			// write to details to journal for this entry in a goroutine
		}
		select {
		case <-agent.shutdown:
			return
		case <-ticker.C:
			//	fmt.Println("Calling Gather .....", time.Now())
			continue
		}
	}
}

func gather(file string, timeout time.Duration, shutdown chan struct{}) {
	timeOutTicker := time.NewTicker(timeout)
	defer timeOutTicker.Stop()
	done := make(chan error)
	go func() {
		done <- parse(file)
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
			fmt.Println("timeout called! parsing not yet completed ... waiting ...", time.Now())
			continue
		}
	}
}
