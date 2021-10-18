package mapbox

import "encoding/json"

type Layer struct {
	ID          string           `json:"id" binding:"required"`
	Type        string           `json:"type" binding:"required,oneof=fill line symbol circle heatmap fill-extrusion raster hillshade background sky"`
	Filter      *json.RawMessage `json:"filter,omitempty"`
	Layout      *json.RawMessage `json:"layout,omitempty"`
	MaxZoom     *float64         `json:"maxzoom,omitempty"`
	MinZoom     *float64         `json:"minzoom,omitempty"`
	Metadata    *json.RawMessage `json:"metadata,omitempty"`
	Paint       *json.RawMessage `json:"paint,omitempty"`
	Source      *string          `json:"source,omitempty"`
	SourceLayer string           `json:"source-layer,omitempty"`
}
