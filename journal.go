package logstream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// A Journal is initialized in-memory and persisted
// periodically to disk. Journal contains JournalEntry
// which keep track to file and last read position.
type Journal struct {
	Path    string
	Entries map[uint32]JournalEntry
	lock  *sync.RWMutex
}

func NewJournal(path string, persist Persist, load Load) (*Journal, error) {
	if path = "" {
			return nil, errors.New("logstream: Journal base path not specified!")
	}
}

// add the any kind of persistance, implementation functions are defined in
// file journal_persistance.go
type Persist func(journal *Journal, waitGroup *sync.WaitGroup) error
type Load func(identifier string) (*Journal, error)

const JOURNAL_ID = "logstream.jrnl"

//var journal *Journal

func CreateJournal(basePath string, persist Persist, load Load) (*Journal, chan JournalEntry, chan bool, error) {
	if basePath == "" {
		return nil, nil, nil, errors.New("logstream: Journal base path not specified!")
	}
	var journal *Journal
	journalFile := filepath.Join(basePath, JOURNAL_ID)
	addChan := make(chan JournalEntry, 50)
	sweepChan := make(chan bool)
	errorChan := make(chan error)

	go func() {
		//var journal *Journal
		if _, err := os.Stat(journalFile); os.IsNotExist(err) {
			//new journal
			journal = &Journal{Path: journalFile,
				Entries: make(map[uint32]JournalEntry),
			}
			fmt.Println("logstream: creating new journal at ", basePath)
		} else {
			//load existing journal
			journal, err = loadFromGob(basePath)
			if err != nil {
				errorChan <- err
			}
		}

		var wg sync.WaitGroup
		for {
			var entry JournalEntry
			var sweep bool
			select {
			case entry = <-addChan:
				journal.Entries[entry.Signature] = entry
			case sweep = <-sweepChan:
				//sweep only if journal has some fresh entries
				if sweep && len(journal.Entries) > 0 {
					wg.Add(1)
					go persist(journal, &wg)
					wg.Wait()
					journal.Entries = make(map[uint32]JournalEntry)
				}
			}
		}
	}()
	return journal, addChan, sweepChan, nil
}

type JournalEntry struct {
	Signature   uint32    `json:"file_signature"`
	File        string    `json:"file_name"`
	Size        int64     `json:"file_size"`
	ModAt       time.Time `json:"file_modified_at"`
	Byte_Offset int64     `json:"offset"`
	Hash        string    `json:"last_hash"`
}

var LOG_LENGTH_FOR_HASH = 500

func NewJournalEntry(signature uint32, file string, size int64, modAt time.Time, offset int64, hash string) *JournalEntry {
	h := make([]byte, LOG_LENGTH_FOR_HASH)
	copy(h, hash[:])
	return &JournalEntry{File: file, Byte_Offset: offset, Hash: string(h)}
}
