package logstream

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateJournal(t *testing.T) {
	_, err := CreateJournal("")
	if err == nil {
		t.Error("logstream: expecting error when location is empty.")
	}
	f, err := CreateJournal("testdir")
	if err != nil {
		t.Error("logstream: expecting a succesful journal creation.")
	}
	if f == nil {
		t.Error("logstream: expecting a journal file.")
	}
	if f.path != "testdir/lgs_jrnl.json" {
		t.Error("logstream: expecting a journal file testdir/lgs_jrnl.json.")
	}
}

func TestJournalWrite(t *testing.T) {
	journal, err := CreateJournal("testdir")
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
		journal.Channel <- je
	}

	time.Sleep(time.Second)
	fmt.Println("Entry 2:", journal.entries)
	if len(journal.entries) != 10 {
			t.Error("logstream: expected 10 journal entries, but found ", len(journal.entries))
	}
}
