package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/davecgh/go-spew/spew"
	"go.i3wm.org/i3"
)

func main() {
	fmt.Println("Hello, World!")

	x, err := xgb.NewConn()
	fatal(err)
	fatal(randr.Init(x))

	heads := getHeads(x)
	_ = heads

	//outputs, _ := i3.GetOutputs()
	//spew.Dump(outputs)

	rcvr := i3.Subscribe(i3.WorkspaceEventType, i3.OutputEventType)
	defer rcvr.Close()

	for rcvr.Next() {
		ev := rcvr.Event()

		switch ev.(type) {
		case *i3.OutputEvent:
			fmt.Printf("Output event: %v\n", ev)
			spew.Dump(ev)

			outputs, _ := i3.GetOutputs()
			spew.Dump(outputs)
		}
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getHeads(x *xgb.Conn) error {
	root := xproto.Setup(x).DefaultScreen(x).Root

	resources, err := randr.GetScreenResourcesCurrent(x, root).Reply()
	if err != nil {
		return err
	}

	//sinfo, _ := randr.GetScreenInfo(x, root).Reply()
	//spew.Dump(sinfo)

	for _, o := range resources.Outputs {
		oinfo, _ := randr.GetOutputInfo(x, o, 0).Reply()

		// Disconnected
		if oinfo.Connection != randr.ConnectionConnected {
			continue
		}

		spew.Dump(oinfo)

		// Off, but how do we get modes for the screens that are off?
		if oinfo.Crtc == 0 {
			continue
		}

		crtc, _ := randr.GetCrtcInfo(x, oinfo.Crtc, 0).Reply()

		spew.Dump(crtc)
		fmt.Println()
		fmt.Println()
	}

	return nil
}
