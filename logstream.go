package logstream

import (
	"bufio"
	"fmt"
	"io"
	"os"
)


type LogFiles struct {
	files []string
}

func (l *LogFiles) Read() (err error) {

	for _, f := range l.files {
		fmt.Println("file :", f)
		fd, err := os.Open(f)
		if err != nil {
			panic(err)
		}
		// close fi on exit and check for its returned error
		defer func() {
			if err := fd.Close(); err != nil {
				panic(err)
			}
		}()

		//var pos int64
    fd.Seek(int64(14),0)
		r := bufio.NewReader(fd)
		for {
			data, err := r.ReadBytes('\n')
			fmt.Println("line :", string(data))
      //pos += int64(len(data))
      //fmt.Println("line 1: ", pos)
      break

			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
		}
	}
  return nil
}

/*func main() {
	logs := LogFiles{files: []string{"/Users/anaray/dev/github/logstream/src/logstream/testdir/logs/l.log"}}
	logs.Read()

}*/
