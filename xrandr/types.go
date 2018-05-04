package xrandr

// QueryOutput represents the output of running `xrandr --query`.
type QueryOutput struct {
	Screens []Screen // TODO(seeruk): Is this hierarchy correct? Can you have multiple screens?
}

// Screen represents an X screen.
type Screen struct {
	Index             int
	MinimumResolution Resolution
	CurrentResolution Resolution
	MaximumResolution Resolution
	Outputs           []Output
}

// Output represents an individual output available to X.
type Output struct {
	Name        string
	IsConnected bool
	Resolution  Resolution
	Position    Position
	Rotation    string // TODO(seeruk): Constant.
	Reflection  string // TODO(seeruk): Constant.
	Modes       []OutputMode
}

// OutputMode represents an output mode. A mode is a resolution and it's supported refresh rates.
type OutputMode struct {
	// Resolution is a resolution supported by an output.
	Resolution Resolution
	// Rates is all of the supported rates for the resolution of this mode. For each resolution
	// there should be one or more rates.
	Rates []Rate
}

// Resolution represents the pixel dimensions of an output, or screen.
type Resolution struct {
	Width  uint
	Height uint
}

// Position represents the pixel position of an output.
type Position struct {
	OffsetX int
	OffsetY int
}

// Rate is a representation of a refresh rate.
type Rate struct {
	// Rate is the actual refresh rate. Sometimes oddly precise (e.g. 59.96).
	Rate float64
	// IsCurrent is true if this rate is currently active on the associated output.
	IsCurrent bool // Represented by `*` in output.
	// IsPreferred is true if this rate is the preferred rate for the associated output.
	IsPreferred bool // Represented by `+` in output.
}

// DP-2 connected 1080x1920+2560+0 left X axis (normal left inverted right x axis y axis) 477mm x 268mm
// Name, IsConnected, Resolution, Position, Rotation, Reflection, Available, Dimensions
