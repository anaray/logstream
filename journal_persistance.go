package logstream

import (
	"fmt"
	"sync"
)

func writeToJson(journal *Journal, wg *sync.WaitGroup) {
	defer wg.Done()

}

func writeToGob(journal *Journal, wg *sync.WaitGroup) {
	defer wg.Done()

}

func writeToStdOut(journal *Journal, wg *sync.WaitGroup) {
	defer wg.Done()
	for key, value := range journal.entries {
		fmt.Println("Key:", key, "Value:", value)
	}
}
