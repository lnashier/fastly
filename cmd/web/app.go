package main

import (
	"fmt"
	"github.com/fastly/cmd/boot"
)

func main() {
	fmt.Println("main enter")
	defer fmt.Println("main exit")
	boot.Up()
}
