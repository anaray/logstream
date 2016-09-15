package logstream

import (
	"os"
	"testing"
)

func TestGetFiles(t *testing.T) {
	files, err := getFiles("", "opendirectoryd.log*")
	if err == nil {
		t.Error("logstream: expecting error")
	}

	files, err = getFiles("testdir/logs/", "opendirectoryd.log*")
	if err != nil {
		t.Error("logstream: expecting non-empty array")
	}
	if files[0] != "testdir/logs/opendirectoryd.log" && files[1] != "testdir/logs/opendirectoryd.log" {
		t.Error("logstream: expecting opendirectoryd.log in array")
	}
}

func TestGenerateFileSignature(t *testing.T) {
	f, _ := os.Open("testdir/logs/sdsc-http.txt")
	signature, _ := getFileSignature(f)
	if signature != 2782789897 {
		t.Error("logstream: expected signature 2782789897 but received", signature)
	}
}

func TestGetFileMetaInfo(t *testing.T) {

}
