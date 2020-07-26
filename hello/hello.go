package main

import (
	"fmt"
)

const (
	spanish = "Spanish"
	french  = "French"
)

const (
	helloPrefixEnglish = "Hello, "
	helloPrefixSpanish = "Hola, "
	helloPrefixFrench  = "Bonjour, "
)

func Hello(name, language string) string {
	if name == "" {
		name = "world"
	}

	prefix := helloPrefixEnglish
	switch language {
	case spanish:
		prefix = helloPrefixSpanish
	case french:
		prefix = helloPrefixFrench
	}

	return fmt.Sprintf("%s%s!", prefix, name)
}

func main() {
	fmt.Println(Hello("", ""))
}
