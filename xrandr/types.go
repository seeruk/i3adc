package xrandr

const (
	RotationNormal = iota
	RotationLeft
	RotationInverted
	RotationRight
)

type Rotation int

const (
	ReflectionXAxis = iota
	ReflectionYAxis
)

type Reflection int

// PropsOutput represents parts of the output of running `xrandr --props`.
type PropsOutput struct {
	Outputs []Output
}

// Output represents an individual output available to X.
type Output struct {
	Name        string
	IsConnected bool
	IsEnabled   bool // If no resolution is set, not enabled.
	Resolution  Resolution
	Position    Position
	Rotation    Rotation
	Reflection  Reflection
	Modes       []OutputMode
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

type Properties struct {
	EDID string // TODO(seeruk): Or should it be bytes?
	// TODO(seeruk): There's a lot more properties. We just don't care about them?
}

// OutputMode represents an output mode. A mode is a resolution and it's supported refresh rates.
type OutputMode struct {
	// Resolution is a resolution supported by an output.
	Resolution Resolution
	// Rates is all of the supported rates for the resolution of this mode. For each resolution
	// there should be one or more rates.
	Rates []Rate
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
// Name, IsConnected, IsEnabled / Resolution, Position, Rotation, Reflection, Available, Dimensions
