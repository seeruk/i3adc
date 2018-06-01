package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/seeruk/i3adc/xrandr"
)

func main() {
	command := exec.Command("xrandr", "--props")
	output, err := command.CombinedOutput()
	fatal(err)
	parser := xrandr.NewParser()

	props, err := parser.ParseProps(output)
	if err != nil {
		log.Println(err)
	}

	for _, output := range props.Outputs {
		if !output.IsConnected || !output.IsEnabled {
			continue
		}

		fmt.Printf("%s: %dx%d\n", output.Name, output.Resolution.Width, output.Resolution.Height)
		fmt.Printf("EDID: %s\n", output.Properties["EDID"])
		fmt.Println()
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
