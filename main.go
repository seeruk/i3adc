package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/seeruk/i3adc/xrandr"
)

func main() {
	command := exec.Command("xrandr", "--props")
	output, err := command.CombinedOutput()
	fatal(err)

	//lexer := xrandr.NewLexer(output)
	//for {
	//	tok := lexer.Scan()
	//	if tok.Type == xrandr.TokenTypeEOF {
	//		break
	//	}
	//
	//	if tok.Type == xrandr.TokenTypeIllegal {
	//		log.Printf("illegal token found at %d:%d: %q", tok.Line, tok.Position, tok.Literal)
	//		break
	//	}
	//
	//	spew.Dump(tok)
	//}
	//
	//os.Exit(1)

	start := time.Now()

	parser := xrandr.NewParser(true)

	for i := 0; i < 1000; i++ {
		_, err := parser.ParseProps(output)
		if err != nil {
			log.Println(err)
		}
	}

	end := time.Since(start)

	fmt.Println(end)

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
