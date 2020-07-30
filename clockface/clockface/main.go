package main

import (
	"os"
	"time"

	"github.com/akaimononoke/learn-go-with-tests-handson/clockface"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
