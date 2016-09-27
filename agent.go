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
	basePath       string
	filterPattern  string
	regexDelim     string
	journalPath    string
}

func NewAgent(basePath, filterPattern, regexDelim string, interval, timeout time.Duration, journalPath string) *Agent {
	agent := Agent{
		gatherInterval: time.Duration(interval),
		gatherTimeout:  time.Duration(timeout),
		shutdown:       make(chan struct{}),
		basePath:       basePath,
		filterPattern:  filterPattern,
		regexDelim:     regexDelim,
		journalPath:    journalPath,
	}
	return &agent
}

func (agent *Agent) Start() {
	//initialize a journal, with gob to serialize to disk
	//journal, journal_chan, sweep_chan, err := CreateJournal(agent.journalPath, writeToGob, loadFromGob)
	journal, err := GetJournal(agent.journalPath, loadFromGob)
	if err != nil {
		panic(err)
	}
	fmt.Println("Got a Journal:", journal)
	//keep the log file delim marker regex pattern compiled
	delim_regex := regexp.MustCompile(agent.regexDelim)
	ticker := time.NewTicker(agent.gatherInterval)
	defer ticker.Stop()
	for {
		files, _ := getFiles(agent.basePath, agent.filterPattern)
		fmt.Println("Selected Files :", files)
		//timeout is a optimistic way to divide the parsing among the
		//files. It is equally divided among the files
		timeOutInSeconds := agent.gatherTimeout.Seconds() / float64(len(files))
		for _, file := range files {
			/*meta, err := getFileMetaInfo(file)
			if err != nil {
				fmt.Println(err)
				continue
			}
			f := meta.signature
			fmt.Println("File Signature :", f)
			je, ok := journal.Get(f)
			//entry is already there in the journal
			if ok {

			}else{

			}
			fmt.Println("Entry :", je)*/
			fmt.Println("Timeout seconds :", timeOutInSeconds)
			gather(file, delim_regex, time.Duration(timeOutInSeconds)*time.Second, agent.shutdown, journal)
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

func gather(file string, delim *regexp.Regexp, timeout time.Duration, shutdown chan struct{}, journal *Journal) {
	//done := make(chan error)
	done := make(chan int64)
	control := make(chan struct{}, 1)
	var t0 time.Time
	timeOutTicker := time.NewTicker(timeout)
	defer timeOutTicker.Stop()

	go func() {
		t0 = time.Now()
		done <- parse(file, delim, control, journal)
	}()

	for {
		select {
		case loc := <-done:
			//update the file seek pos
			fmt.Println("parsing completed !!!! file:", file, " location: ", loc, " time:", time.Now())
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
			//break
			//return
		}
	}
	return

}
