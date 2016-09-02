package logstream

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func writeToJson(journal *Journal, wg *sync.WaitGroup) error {
	defer wg.Done()
	return nil
}

func writeToGob(journal *Journal, wg *sync.WaitGroup) error {
	defer wg.Done()
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(journal)
	if err != nil {
		return err
	}
	f, eopen := os.OpenFile(journal.Path, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if eopen != nil {
		return eopen
	}
	n, e := f.Write(b.Bytes())
	if e != nil {
		return e
	}
	fmt.Fprintf(os.Stderr, "%d bytes successfully written to file\n", n)
	return nil
}

func loadFromGob(basePath string) (journal *Journal, err error) {
	journalFile := filepath.Join(basePath, JOURNAL_ID)
	f, err := os.Open(journalFile)
	if err != nil {
		return nil, err
	}
	j := new(Journal)
	dec := gob.NewDecoder(f)
	err = dec.Decode(&j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func writeToStdOut(journal *Journal, wg *sync.WaitGroup) error {
	defer wg.Done()
	for key, value := range journal.Entries {
		fmt.Println("Key:", key, "Value:", value)
	}
	return nil
}
