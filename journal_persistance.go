package logstream

import (
	"fmt"
	"sync"
)

func writeToJson(val *map[uint64]JournalEntry, wg *sync.WaitGroup) {
	defer wg.Done()

}

func writeToGob(val *map[uint64]JournalEntry, wg *sync.WaitGroup) {
	defer wg.Done()

}

func writeToStdOut(val *map[uint64]JournalEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	for key, value := range *val {
		fmt.Println("Key:", key, "Value:", value)
	}
}
