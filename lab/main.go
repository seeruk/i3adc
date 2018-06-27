// Example randr uses the randr protocol to get information about the active
// heads. It also listens for events that are sent when the head configuration
// changes. Since it listens to events, you'll have to manually kill this
// process when you're done (i.e., ctrl+c.)
//
// While this program is running, if you use 'xrandr' to reconfigure your
// heads, you should see event information dumped to standard out.
//
// For more information, please see the RandR protocol spec:
// http://www.x.org/releases/X11R7.6/doc/randrproto/randrproto.txt
package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/davecgh/go-spew/spew"
)

type Mode struct {
	Name     string
	ModeInfo randr.ModeInfo
}

func main() {
	X, _ := xgb.NewConn()

	// Every extension must be initialized before it can be used.
	err := randr.Init(X)
	if err != nil {
		log.Fatal(err)
	}

	// Get the root window on the default screen.
	root := xproto.Setup(X).DefaultScreen(X).Root

	// Gets the current screen resources. Screen resources contains a list
	// of names, crtcs, outputs and modes, among other things.
	resources, err := randr.GetScreenResources(X, root).Reply()
	if err != nil {
		log.Fatal(err)
	}

	var nameOffset int

	modes := make(map[uint32]Mode)
	for _, xmode := range resources.Modes {
		modes[xmode.Id] = Mode{
			Name:     string(resources.Names[nameOffset : nameOffset+int(xmode.NameLen)]),
			ModeInfo: xmode,
		}

		nameOffset += int(xmode.NameLen)
	}

	primary, err := randr.GetOutputPrimary(X, root).Reply()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through all of the outputs and show some of their info.
	for _, output := range resources.Outputs {
		info, err := randr.GetOutputInfo(X, output, 0).Reply()
		if err != nil {
			log.Fatal(err)
		}

		var crtcLine string

		crtcInfo, err := randr.GetCrtcInfo(X, info.Crtc, 0).Reply()
		if err == nil {
			// This can fail, if the screen is not active.
			crtcLine = fmt.Sprintf("%dx%d+%d+%d", crtcInfo.Width, crtcInfo.Height, crtcInfo.X, crtcInfo.Y)
		}

		connectionStatus := "connected"
		if info.Connection != randr.ConnectionConnected {
			connectionStatus = "disconnected"
		}

		primaryStatus := ""
		if output == primary.Output {
			primaryStatus = "primary"
		}

		fmt.Printf("%s %s %s %s\n",
			string(info.Name),
			connectionStatus,
			primaryStatus,
			crtcLine,
		)

		props, err := randr.ListOutputProperties(X, output).Reply()
		if err != nil {
			log.Fatal(err)
		}

		for _, atom := range props.Atoms {
			propContent, err := randr.GetOutputProperty(X, output, atom, 0, 0, 0, false, false).Reply()
			if err != nil {
				log.Fatal(err)
			}

			propContent, err = randr.GetOutputProperty(X, output, atom, 0, 0, propContent.BytesAfter, false, false).Reply()
			if err != nil {
				log.Fatal(err)
			}

			atomData, err := xproto.GetAtomName(X, atom).Reply()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("    %s:\n", atomData.Name)
			fmt.Printf("        %x\n", propContent.Data)
		}

		for i, modeID := range info.Modes {
			mode := modes[uint32(modeID)]
			fmt.Printf("    %s: width: %d, height: %d", mode.Name, mode.ModeInfo.Width, mode.ModeInfo.Height)
			if crtcInfo != nil && uint32(crtcInfo.Mode) == mode.ModeInfo.Id {
				fmt.Print(" +current")
			}
			if i < int(info.NumPreferred) {
				fmt.Printf(" +preferred")
			}
			fmt.Print("\n")

			spew.Dump(mode)
		}
	}

	fmt.Println("\n")

	// Iterate through all of the crtcs and show some of their info.
	for _, crtc := range resources.Crtcs {
		info, err := randr.GetCrtcInfo(X, crtc, 0).Reply()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("X: %d, Y: %d, Width: %d, Height: %d\n",
			info.X, info.Y, info.Width, info.Height)

		var rotation string
		var reflection string

		// Rotation:
		switch {
		case (info.Rotation & randr.RotationRotate0) != 0:
			rotation = "normal"
		case (info.Rotation & randr.RotationRotate90) != 0:
			rotation = "left"
		case (info.Rotation & randr.RotationRotate180) != 0:
			rotation = "inverted"
		case (info.Rotation & randr.RotationRotate270) != 0:
			rotation = "right"
		}

		switch {
		case (info.Rotation & randr.RotationReflectX) != 0:
			reflection = "x"
		case (info.Rotation & randr.RotationReflectY) != 0:
			reflection = "y"
		default:
			reflection = "normal"
		}

		fmt.Printf("rotation: %s, reflection %s\n", rotation, reflection)

		for _, output := range info.Outputs {
			info, err := randr.GetOutputInfo(X, output, 0).Reply()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(info.Name))
		}
	}
}
