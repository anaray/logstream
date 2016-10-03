# logstream 
Agent for Log Harvesting 
status: work in progress

Logstream is a light weight agent for parsing and moving the logs from the source to any configured destination. 

```
Usage:
logstream -conf=my_conf.json

configuration json contains the necessary configuration for logstream and configuration parameters are:

{
  "journal_path" : "/my_journal_loc/",
  "base_directory" : "/var/log/",
  "filter_pattern" : "displaypolicyd.log",
  "log_delim_regex" : "u>\\d*",
  "log_type" : "displaypolicyd",
  "interval" : 60
}
```

* journal_path = a file stored in disk to maintain file read position and other meta-data like size, modified_at etc.
* base_directory = directory where files to be read are located.
* filter_pattern = regex file filters.
* log_delim_regex = regular expression pattern to identify each log entries.
* log_type = a textual description of type/kind of the logfile. example access, application etc.
* interval = how often log files are harvested (seconds).

Build:
```
Create a golang project. https://golang.org/doc/code.html
go get github.com/anaray/logstream/
cd $GOPATH/src/logstream/
./build.sh - it creates executable logstream_<OS>_<ARCH>
```
