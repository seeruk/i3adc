package props

const (
	RotationNormal = iota
	RotationLeft
	RotationInverted
	RotationRight
)

type Rotation int

const (
	ReflectionNone = iota
	ReflectionXAxis
	ReflectionYAxis
)

type Reflection int

// CommandOutput represents parts of the output of running `xrandr --props`.
type CommandOutput struct {
	Outputs []Output `json:"outputs"`
}

// Output represents an individual output available to X.
type Output struct {
	Name        string            `json:"name"`
	IsConnected bool              `json:"is_connected"`
	IsPrimary   bool              `json:"is_primary"`
	IsEnabled   bool              `json:"is_enabled"` // If no resolution is set, not enabled.
	Resolution  Resolution        `json:"resolution"`
	Position    Position          `json:"position"`
	Rotation    Rotation          `json:"rotation"`
	Reflection  Reflection        `json:"reflection"`
	Dimensions  Dimensions        `json:"dimensions"`
	Properties  map[string]string `json:"properties"`
	Modes       []OutputMode      `json:"modes"`
}

// Resolution represents the pixel dimensions of an output, or screen.
type Resolution struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

// Position represents the pixel position of an output.
type Position struct {
	OffsetX int `json:"offset_x"`
	OffsetY int `json:"offset_y"`
}

// Dimensions represents the physical dimensions of an output. Usually in mm.
type Dimensions struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

// OutputMode represents an output mode. A mode is a resolution and it's supported refresh rates.
type OutputMode struct {
	// Resolution is a resolution supported by an output.
	Resolution Resolution `json:"resolution"`
	// Rates is all of the supported rates for the resolution of this mode. For each resolution
	// there should be one or more rates.
	Rates []Rate `json:"rates"`
}

// Rate is a representation of a refresh rate.
type Rate struct {
	// Rate is the actual refresh rate. Sometimes oddly precise (e.g. 59.96).
	Rate float64 `json:"rate"`
	// IsCurrent is true if this rate is currently active on the associated output.
	IsCurrent bool `json:"is_current"` // Represented by `*` in output.
	// IsPreferred is true if this rate is the preferred rate for the associated output.
	IsPreferred bool `json:"is_preferred"` // Represented by `+` in output.
}
