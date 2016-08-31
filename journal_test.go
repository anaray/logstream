package logstream

import (
	"testing"
	"time"
	"sync"
)

func TestCreateJournal(t *testing.T) {
	_, _, err := CreateJournal("", writeToStdOut)
	if err == nil {
		t.Error("logstream: expecting error when location is empty.")
	}

	a, s, err := CreateJournal("testdir", writeToStdOut)
	if a == nil || s == nil {
		t.Error("logstream: expecting 2 channels, one for adding Journal entry and one for sweep event.")
	}
	if err != nil {
		t.Error("logstream: expecting a succesful journal creation.")
	}
}

func TestJournalWrite(t *testing.T) {

	f := func(journal *Journal, wg *sync.WaitGroup){
		defer wg.Done()
		if len(journal.entries) != 10 {
			t.Error("logstream: expecting 10 journal entries.")
			}
	}
	a, s, err := CreateJournal("testdir", f)
	if err != nil {
		t.Error("logstream: expecting a succesful journal creation.")
	}
	var je JournalEntry
	var i uint64
	for i = 0; i < 10; i++ {
		je = JournalEntry{}
		je.Ino = i
		je.Byte_Offset = 10
		je.Hash = "1FDAA"
		a <- je
	}
	time.Sleep(time.Second)
	s <- true
	time.Sleep(time.Second)
}
