package main

import (
	"os"
	"time"

	"github.com/akaimononoke/learn-go-with-tests-handson/foundation/math"
)

func main() {
	t := time.Now()
	math.SVGWriter(os.Stdout, t)
}
