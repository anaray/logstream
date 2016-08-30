package logstream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Register struct {
	journals []*Journal
}

// A Journal is initialized in-memory and persisted
// periodically to disk. Journal contains JournalEntry
// which keep track to file and last read position.
type Journal struct {
	path    string
	entries map[uint64]JournalEntry
}

// add the any kind of persistance, implementation functions are defined in
// file journal_persistance.go
type Persist func(journalEntries *map[uint64]JournalEntry, waitGroup *sync.WaitGroup)

var JOURNAL_ID = "lgs_jrnl.json"

func CreateJournal(basePath string, persist Persist) (chan JournalEntry, chan bool, error) { //(*Journal, error) {
	if basePath == "" {
		return nil, nil, errors.New("logstream: Journal base path not specified!")
	}
	journalFile := filepath.Join(basePath, JOURNAL_ID)
	file, file_err := os.OpenFile(journalFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if file_err != nil {
		err_str := fmt.Errorf("logstream: Error opening logstream journal path: %s", file_err.Error())
		return nil, nil, errors.New(err_str.Error())
	}
	defer file.Close()
	addChan := make(chan JournalEntry, 50)
	sweepChan := make(chan bool)

	go func() {
		//create the in-memory journal
		journal := Journal{path: journalFile,
			entries: make(map[uint64]JournalEntry),
		}
		var wg sync.WaitGroup
		for {
			var entry JournalEntry
			var sweep bool
			select {
			case entry = <-addChan:
				journal.entries[entry.Ino] = entry
			case sweep = <-sweepChan:
				if sweep {
					wg.Add(1)
					go persist(&journal.entries, &wg)
				}
			}
		}
	}()
	return addChan, sweepChan, nil
}

/*func (journal *Journal) Write(entry *JournalEntry) error {
	if entry != nil {
		e, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		file, file_err := os.OpenFile(journal.path, os.O_APPEND|os.O_WRONLY, 0660)
		if file_err != nil {
			err_str := fmt.Errorf("logstream: Error opening logstream journal path: %s", file_err.Error())
			return errors.New(err_str.Error())
		}
		defer file.Close()
		if _, err = file.Write(append(e, '\n')); err != nil {
			return err
		}
	}
	return nil
}*/

//add major & minor number
type JournalEntry struct {
	Ino         uint64 `json:"inode"`
	Byte_Offset int64  `json:"offset"`
	Hash        string `json:"last_hash"`
}

var LOG_LENGTH_FOR_HASH = 500

func NewJournalEntry(inode uint64, offset int64, hash string) *JournalEntry {
	h := make([]byte, LOG_LENGTH_FOR_HASH)
	copy(h, hash[:])
	return &JournalEntry{Ino: inode, Byte_Offset: offset, Hash: string(h)}
}
