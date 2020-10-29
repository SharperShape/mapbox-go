package mapbox

// Transition controls timing for the interpolation between a style property's previous value
// and new value. A style's root transition property provides global transition defaults for that style.
type Transition struct {
	// Duration is the time allotted for transitions to complete.
	Duration int64
	// Delay is the length of time before a transition begins.
	Delay int64
}
