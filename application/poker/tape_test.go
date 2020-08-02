package poker_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/akaimononoke/learn-go-with-tests-handson/application/poker"
)

func createTempFile(t *testing.T, data string) (*os.File, func()) {
	t.Helper()

	tmp, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	tmp.WriteString(data)

	removeFile := func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}

	return tmp, removeFile
}

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()
	tape := &poker.Tape{file}

	want := "abc"

	tape.Write([]byte(want))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)
	got := string(newFileContents)

	if want != got {
		t.Errorf("written file content is invalid: want %q, got %q", want, got)
	}
}
