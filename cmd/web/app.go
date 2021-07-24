package main

import (
	"fmt"
	"github.com/fastly/cmd/boot"
	"github.com/spf13/viper"
	"os"
)

func main() {
	fmt.Println("main enter")
	defer fmt.Println("main exit")

	var cfg *viper.Viper
	var err error
	if cfg, err = boot.Setup(); err != nil {
		fmt.Printf("main failed to setup config %s\n", err.Error())
		os.Exit(1)
	}

	boot.Up(cfg)
}
