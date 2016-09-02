package logstream

import (
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestCreateJournal(t *testing.T) {
	dummyLoadFn := func(file string) (*Journal, error) {
		return nil, nil
	}
	_, _, err := CreateJournal("", writeToStdOut, dummyLoadFn)
	if err == nil {
		t.Error("logstream: expecting error when location is empty.")
	}

	a, s, err := CreateJournal("testdir", writeToStdOut, dummyLoadFn)
	if a == nil || s == nil {
		t.Error("logstream: expecting 2 channels, one for adding Journal entry and one for sweep event.")
	}
	if err != nil {
		t.Error("logstream: expecting a succesful journal creation.")
	}
}

func TestJournalWrite(t *testing.T) {
	dummyLoadFn := func(file string) (*Journal, error) {
		return nil, nil
	}
	writeFn := func(journal *Journal, wg *sync.WaitGroup) error {
		defer wg.Done()
		if len(journal.Entries) != 10 {
			t.Error("logstream: expecting 10 journal entries.")
		}
		return nil
	}
	a, s, err := CreateJournal("testdir", writeFn, dummyLoadFn)
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

func TestJournalPersistance(t *testing.T) {
	a, s, err := CreateJournal("testdir", writeToGob, loadFromGob)
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

func TestWriteAndLoadToAndFromGob(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	journalFile := filepath.Join("testdir", JOURNAL_ID)
	journal := Journal{Path: journalFile,
		Entries: make(map[uint64]JournalEntry),
	}
	je := JournalEntry{}
	je.Ino = 10
	je.Byte_Offset = 10
	je.Hash = "1FDAA"

	journal.Entries[1] = je
	err := writeToGob(&journal, &wg)

	if err != nil {
		t.Error("logstream: write to journal failed.")
	}
	j, err := loadFromGob("testdir")
	e := j.Entries[1]
	if err != nil {
		t.Error("logstream: loading journal failed.")
	}
	if e.Ino != 10 {
		t.Error("logstream: loading journal failed, expected Ino value 10.")
	}
	if e.Byte_Offset != 10 {
		t.Error("logstream: loading journal failed, expected Byte_Offset value 10.")
	}
	if e.Hash != "1FDAA" {
		t.Error("logstream: loading journal failed, expected Hash value 1FDAA.")
	}
}
