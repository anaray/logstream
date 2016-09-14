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
	Interval      int64  `json:"interval"`
	Timeout       int64  `json:"timeout"`
}

func main() {
	f := flag.String("conf", ".", "the configuration file")
	flag.Parse()
	c, err := ioutil.ReadFile(*f)
	if err != nil {
		panic("couldn't read file.")
		os.Exit(1)
	}
	conf := Conf{}
	json.Unmarshal(c, &conf)
	fmt.Println(conf)

	agent := logstream.NewAgent(conf.Path,
		conf.FilterPattern,
		time.Duration(conf.Interval)*time.Second,
		time.Duration(conf.Timeout)*time.Second)
	agent.Start()
}
