package poker

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()
	tape := &tape{file}

	want := "abc"

	tape.Write([]byte(want))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)
	got := string(newFileContents)

	if want != got {
		t.Errorf("written file content is invalid: want %q, got %q", want, got)
	}
}
