package main

import (
	"fmt"

	"go.i3wm.org/i3"
)

func main() {
	fmt.Println("Hello, World!")

	rcvr := i3.Subscribe(i3.WorkspaceEventType, i3.OutputEventType)
	defer rcvr.Close()

	for rcvr.Next() {
		ev := rcvr.Event()
		//spew.Dump(ev)

		switch ev.(type) {
		case *i3.WorkspaceEvent:
			fmt.Printf("Workspace event: %v\n", ev)
		case *i3.OutputEvent:
			fmt.Printf("Output event: %v\n", ev)
		}
	}
}
