package maps

import (
	"testing"
)

func TestSearch(t *testing.T) {
	dict := Dictionary{"test": "this is just a test"}

	t.Parallel()

	t.Run("known word", func(t *testing.T) {
		got, _ := dict.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, err := dict.Search("unknown")
		assertError(t, err, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	t.Parallel()

	t.Run("new word", func(t *testing.T) {
		dict := Dictionary{}
		word, definition := "test", "this is just a test"
		err := dict.Add(word, definition)

		assertError(t, err, nil)
		assertDefinition(t, dict, word, definition)
	})

	t.Run("existing word", func(t *testing.T) {
		word, definition := "test", "this is just a test"
		dict := Dictionary{word: definition}
		err := dict.Add(word, "new test")

		assertError(t, err, ErrWordExists)
		assertDefinition(t, dict, word, definition)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("existing word", func(t *testing.T) {
		word, definition := "test", "this is just a test"
		dict := Dictionary{word: definition}
		newDefinition := "new definition"

		err := dict.Update(word, newDefinition)

		assertError(t, err, nil)
		assertDefinition(t, dict, word, newDefinition)
	})

	t.Run("new word", func(t *testing.T) {
		word, definition := "test", "this is just a test"
		dict := Dictionary{}

		err := dict.Update(word, definition)

		assertError(t, err, ErrWordDoesNotExist)
	})
}

func TestDelete(t *testing.T) {
	word := "test"
	dict := Dictionary{word: "test definition"}

	dict.Delete(word)

	_, err := dict.Search(word)
	if err != ErrNotFound {
		t.Errorf("Expected %q to be deleted", word)
	}
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q, want %q", got, want)
	}
	if got == nil && want != nil {
		t.Fatal("expected to get an error")
	}
}

func assertDefinition(t *testing.T, dict Dictionary, word, definition string) {
	t.Helper()

	got, err := dict.Search(word)
	if err != nil {
		t.Fatalf("should find added word: %v", err)
	}
	if got != definition {
		t.Errorf("got %q, want %q", got, definition)
	}
}
