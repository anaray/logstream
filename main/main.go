package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/anaray/logstream/agent"
	"io/ioutil"
	"os"
	"time"
	//"runtime"
)

type Conf struct {
	Path          string `toml:"base_directory"`
	FilterPattern string `toml:"filter_pattern"`
	LogDelimRegex string `toml:"log_delim_regex"`
	LogType       string `toml:"log_type"`
	Interval      int64  `toml:"interval"`
	Timeout       int64  `toml:"timeout"`
	JournalPath   string `toml:"journal_path"`
}

/*func init() {
  runtime.GOMAXPROCS(runtime.NumCPU())
}*/

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
	if _, err := toml.Decode(string(c), &conf); err != nil {
		panic(fmt.Sprintf("logstream: failed to parse configuration file %s\n", *f))
		os.Exit(1)
	}
	logger.Logf("Creating LogStream Agent ...")
	logger.Logf("File Path: %s\n", conf.Path)
	logger.Logf("Filter: %s\n", conf.FilterPattern)
	logger.Logf("Log Type: %s\n", conf.LogType)
	logger.Logf("Log Delim Regex: %s\n", conf.LogDelimRegex)
	logger.Logf("Gather Interval: %d\n", conf.Interval)
	logger.Logf("Gather Timeout: %d\n", conf.Timeout)
	logger.Logf("Journal Path: %s\n", conf.JournalPath)
	agent := logstream.NewAgent(conf.Path,
		conf.FilterPattern,
		conf.LogDelimRegex,
		conf.LogType,
		time.Duration(conf.Interval)*time.Second,
		time.Duration(conf.Timeout)*time.Second,
		conf.JournalPath)
	go func() {
		for {
			r := <-agent.Output
            //push to downstream plugin
            //TODO: implement output plugin
			fmt.Println("Log entry :", r)
		}
	}()
	agent.Start()
}
