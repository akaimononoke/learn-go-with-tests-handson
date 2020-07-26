package main

import (
	"fmt"
)

const helloPrefixEn = "Hello, "

func Hello(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("%s%s!", helloPrefixEn, name)
}

func main() {
	fmt.Println(Hello(""))
}
