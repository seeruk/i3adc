package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	fmt.Println("Hello, World!")

	command := exec.Command("xrandr", "--query")
	output, err := command.CombinedOutput()
	fatal(err)

	fmt.Println(string(output))
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
