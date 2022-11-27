package main

import (
	"fmt"
	"os"

	"github.com/alinnert/xt/commands"
)

func main() {
	mainCommand := commands.MainCommand()

	if err := mainCommand.Execute(); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
