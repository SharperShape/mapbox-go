package mapbox

// https://docs.mapbox.com/mapbox-gl-js/style-spec/light/
type Light struct {
	Anchor    string  `json:"anchor,omitempty"`
	Color     string  `json:"color,omitempty"`
	Intensity float64 `json:"intensity,omitempty"`
}
