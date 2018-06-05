package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/i3adc/xrandr"
)

func main() {
	command := exec.Command("xrandr", "--props")
	output, err := command.CombinedOutput()
	fatal(err)

	lexer := xrandr.NewLexer(output)
	for {
		tok := lexer.Scan()
		if tok.Type == xrandr.TokenTypeEOF {
			break
		}

		if tok.Type == xrandr.TokenTypeIllegal {
			log.Printf("illegal token found at %d:%d: %q", tok.Line, tok.Position, tok.Literal)
			break
		}

		spew.Dump(tok)
	}

	os.Exit(1)

	//parser := xrandr.NewParser()
	//
	//props, err := parser.ParseProps(output)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//for _, output := range props.Outputs {
	//	if !output.IsConnected || !output.IsEnabled {
	//		continue
	//	}
	//
	//	fmt.Printf("%s: %dx%d\n", output.Name, output.Resolution.Width, output.Resolution.Height)
	//	fmt.Printf("EDID: %s\n", output.Properties["EDID"])
	//	fmt.Println()
	//}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
