package main

import (
	"fmt"
)

const (
	spanish  = "Spanish"
	french   = "French"
	japanese = "Japanese"
)

func Hello(name, language string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf(helloFormatByLanguage(language), name)
}

func helloFormatByLanguage(language string) (format string) {
	switch language {
	case spanish:
		format = "Hola, %s!"
	case french:
		format = "Bonjour, %s!"
	case japanese:
		format = "こんにちは、%s！"
	default:
		format = "Hello, %s!"
	}
	return
}

func main() {
	fmt.Println(Hello("", ""))
}
