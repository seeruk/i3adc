package xrandr

// Rotation values.
const (
	RotationNormal Rotation = iota
	RotationLeft
	RotationInverted
	RotationRight
)

// Rotation represents the rotation status of a CRTC, which is sent to an output.
type Rotation int

// Reflection values.
const (
	ReflectionNormal Reflection = iota
	ReflectionX
	ReflectionY
)

// Reflection represents the reflection status of a CRTC, which is sent to an output.
type Reflection int

// Output represents an randr output, condensing the information we need into a simple struct.
type Output struct {
	Name        string     `json:"name"`
	IsConnected bool       `json:"is_connected"`
	IsEnabled   bool       `json:"is_enabled"`
	IsPrimary   bool       `json:"is_primary"`
	ModeName    string     `json:"mode_name"`
	Width       uint       `json:"width_px"`
	Height      uint       `json:"height_px"`
	OffsetX     int        `json:"offset_x"`
	OffsetY     int        `json:"offset_y"`
	Rotation    Rotation   `json:"rotation"`
	Reflection  Reflection `json:"reflection"`
	Properties  Properties `json:"properties,omitempty"`
	Modes       []Mode     `json:"modes"`
}

// Properties represents the properties of an output.
type Properties map[string][]byte

// Mode represents an randr mode, only including the information we need.
type Mode struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Width       uint   `json:"width"`
	Height      uint   `json:"height"`
	IsPreferred bool   `json:"is_preferred"`
}
