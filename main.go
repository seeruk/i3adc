package main

import (
	"fmt"
	"log"

	"github.com/seeruk/i3adc/xrandr"
)

func main() {
	hash, err := xrandr.All()
	fatal(err)

	fmt.Println(hash)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
