package main

import (
	"bytes"
	"reflect"
	"testing"
)

type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

type Call int

const (
	write Call = iota
	sleep
)

type SpyWriteSleeper struct{ Calls []Call }

func (s *SpyWriteSleeper) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *SpyWriteSleeper) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

func TestCountDown(t *testing.T) {
	t.Parallel()

	t.Run("prints 3 to Go!", func(t *testing.T) {
		buf := &bytes.Buffer{}
		sleeper := &SpySleeper{}
		CountDown(buf, sleeper)

		got := buf.String()
		want := `3
2
1
Go!`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("sleep before every print", func(t *testing.T) {
		writeSleeper := &SpyWriteSleeper{}
		CountDown(writeSleeper, writeSleeper)

		want := []Call{sleep, write, sleep, write, sleep, write, sleep, write}
		if !reflect.DeepEqual(writeSleeper.Calls, want) {
			t.Errorf("got %#v, want %#v", writeSleeper.Calls, want)
		}
	})
}
