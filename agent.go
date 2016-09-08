package logstream

import (
	"fmt"
	"time"
)

type Agent struct {
	gatherInterval time.Duration
	gatherTimeout  time.Duration
	shutdown       chan struct{}
}

func NewAgent(interval, timeout time.Duration) *Agent {
	agent := Agent{
		gatherInterval: time.Duration(interval),
		gatherTimeout:  time.Duration(timeout),
		shutdown:       make(chan struct{}),
	}
	return &agent
}

func (agent *Agent) Start(file string, seekPos int64) {
	ticker := time.NewTicker(agent.gatherInterval)
	defer ticker.Stop()
	for {
		gather(file, agent.gatherTimeout, agent.shutdown)
		select {
		case <-agent.shutdown:
			return
		case <-ticker.C:
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
			fmt.Println("parsing completed !!!!")
			return
		case <-shutdown:
			fmt.Println("shutdown called !")
			return
		case <-timeOutTicker.C:
			fmt.Println("timeout called! parsing not yet completed ... waiting ...")
			continue
		}
	}
}
