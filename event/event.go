package event

// Event represents an i3adc event.
type Event struct {
	// IsStartup is used to distinguish between the startup event, and regular events.
	IsStartup bool
}
