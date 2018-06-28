package xrandr

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

// Client is an xrandr client, using the X binary protocol for communication via XGB.
type Client struct {
	conn *xgb.Conn
	root xproto.Window
}

// NewClient returns a new instance of the xrandr client.
func NewClient() (*Client, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	// Every extension must be initialized before it can be used.
	err = randr.Init(conn)
	if err != nil {
		return nil, err
	}

	// Get the root window on the default screen.
	root := xproto.Setup(conn).DefaultScreen(conn).Root

	return &Client{
		conn: conn,
		root: root,
	}, nil
}

func (c *Client) GetOutputs() ([]Output, error) {
	var outputs []Output

	// Gets the current screen resources. Screen resources contains a list
	// of names, crtcs, outputs and modes, among other things.
	resources, err := randr.GetScreenResources(c.conn, c.root).Reply()
	if err != nil {
		return outputs, err
	}

	primary, err := randr.GetOutputPrimary(c.conn, c.root).Reply()
	if err != nil {
		return outputs, err
	}

	modes := prepareModes(resources)

	// Iterate through all of the outputs and show some of their info.
	for _, xoutput := range resources.Outputs {
		info, err := randr.GetOutputInfo(c.conn, xoutput, 0).Reply()
		if err != nil {
			return outputs, err
		}

		output := Output{}
		output.Name = string(info.Name)
		output.Properties = make(Properties)

		if info.Connection == randr.ConnectionConnected {
			output.IsConnected = true
		}

		if xoutput == primary.Output {
			output.IsPrimary = true
		}

		for i, modeID := range info.Modes {
			mode := modes[uint32(modeID)]
			if i < int(info.NumPreferred) {
				mode.IsPreferred = true
			}

			output.Modes = append(output.Modes, mode)
		}

		// Get information about the currently active mode, apply it to the output if possible.
		crtcInfo, err := randr.GetCrtcInfo(c.conn, info.Crtc, 0).Reply()
		if err == nil {
			// It's okay for this one to error.
			output.IsEnabled = true
			output.Width = uint(crtcInfo.Width)
			output.Height = uint(crtcInfo.Height)
			output.OffsetX = int(crtcInfo.X)
			output.OffsetY = int(crtcInfo.Y)

			// Assign mode name, so we can set it again later via xrandr command.
			for _, modeID := range info.Modes {
				mode := modes[uint32(modeID)]
				if uint(crtcInfo.Mode) == mode.ID {
					output.ModeName = mode.Name
				}
			}

			// Rotation:
			switch {
			case (crtcInfo.Rotation & randr.RotationRotate0) != 0:
				output.Rotation = RotationNormal
			case (crtcInfo.Rotation & randr.RotationRotate90) != 0:
				output.Rotation = RotationLeft
			case (crtcInfo.Rotation & randr.RotationRotate180) != 0:
				output.Rotation = RotationInverted
			case (crtcInfo.Rotation & randr.RotationRotate270) != 0:
				output.Rotation = RotationRight
			}

			// Reflection:
			switch {
			case (crtcInfo.Rotation & randr.RotationReflectX) != 0:
				output.Reflection = ReflectionX
			case (crtcInfo.Rotation & randr.RotationReflectY) != 0:
				output.Reflection = ReflectionY
			default:
				output.Reflection = ReflectionNormal
			}
		}

		// Assign properties to the output.
		props, err := randr.ListOutputProperties(c.conn, xoutput).Reply()
		if err != nil {
			return outputs, err
		}

		for _, atom := range props.Atoms {
			// Fetch the property once, with no bytes so we can see how many bytes we need to request.
			propContent, err := randr.GetOutputProperty(c.conn, xoutput, atom, 0, 0, 0, false, false).Reply()
			if err != nil {
				return outputs, err
			}

			// Then fetch the property, asking for that many bytes of data.
			propContent, err = randr.GetOutputProperty(c.conn, xoutput, atom, 0, 0, propContent.BytesAfter, false, false).Reply()
			if err != nil {
				return outputs, err
			}

			atomData, err := xproto.GetAtomName(c.conn, atom).Reply()
			if err != nil {
				return outputs, err
			}

			output.Properties[atomData.Name] = propContent.Data
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

// prepareModes gets all of the modes for the current screen in a format that's more useful to us.
func prepareModes(resources *randr.GetScreenResourcesReply) map[uint32]Mode {
	var nameOffset int

	modes := make(map[uint32]Mode)
	for _, xmode := range resources.Modes {
		modes[xmode.Id] = Mode{
			ID:          uint(xmode.Id),
			Name:        string(resources.Names[nameOffset : nameOffset+int(xmode.NameLen)]),
			Width:       uint(xmode.Width),
			Height:      uint(xmode.Height),
			IsPreferred: false,
		}

		nameOffset += int(xmode.NameLen)
	}

	return modes
}
