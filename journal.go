package logstream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// A Journal is initialized in-memory and persisted
// periodically to disk. Journal contains JournalEntry
// which keep track to file and last read position.
type Journal struct {
	Path    string
	Entries map[string]JournalEntry
}

// add the any kind of persistance, implementation functions are defined in
// file journal_persistance.go
type Persist func(journal *Journal, waitGroup *sync.WaitGroup) error
type Load func(identifier string) (*Journal, error)

const JOURNAL_ID = "logstream.jrnl"

var journal *Journal

func CreateJournal(basePath string, persist Persist, load Load) (chan JournalEntry, chan bool, error) {
	if basePath == "" {
		return nil, nil, errors.New("logstream: Journal base path not specified!")
	}
	journalFile := filepath.Join(basePath, JOURNAL_ID)
	addChan := make(chan JournalEntry, 50)
	sweepChan := make(chan bool)

	go func() {
		//var journal *Journal
		if _, err := os.Stat(journalFile); os.IsNotExist(err) {
			//new journal
			journal = &Journal{Path: journalFile,
				Entries: make(map[string]JournalEntry),
			}
			fmt.Println("logstream: creating new journal at ", basePath)
		} else {
			//load existing journal
			journal, err = loadFromGob(basePath)
		}

		var wg sync.WaitGroup
		for {
			var entry JournalEntry
			var sweep bool
			select {
			case entry = <-addChan:
				journal.Entries[entry.File] = entry
			case sweep = <-sweepChan:
				if sweep {
					wg.Add(1)
					go persist(journal, &wg)
					wg.Wait()
					journal.Entries = make(map[string]JournalEntry)
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
	//Ino         uint64 `json:"inode"`
	File	string `json:"file"`
	Byte_Offset int64  `json:"offset"`
	Hash        string `json:"last_hash"`
}

var LOG_LENGTH_FOR_HASH = 500

func NewJournalEntry(file string, offset int64, hash string) *JournalEntry {
	h := make([]byte, LOG_LENGTH_FOR_HASH)
	copy(h, hash[:])
	return &JournalEntry{File: file, Byte_Offset: offset, Hash: string(h)}
}
