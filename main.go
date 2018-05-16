package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/i3adc/xrandr"
)

func main() {
	fmt.Println("Hello, World!")

	command := exec.Command("xrandr", "--props")
	output, err := command.CombinedOutput()
	fatal(err)

	lexer := xrandr.NewLexer(output)

	for {
		tok, err := lexer.Scan()
		if err != nil {
			log.Println(err)
			break
		}

		spew.Dump(tok)

		if tok.Type == xrandr.TokenTypeEOF {
			break
		}
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
