package logstream

import (
	"testing"
	"time"
)

func TestAgent(t *testing.T) {
	agent := NewAgent(time.Minute, time.Duration(30)*time.Second)
	agent.Start("testdir/logs/opendirectoryd.log", 0)
}
