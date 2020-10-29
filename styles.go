package mapbox

import (
	"encoding/json"
	"path"
	"time"

	"github.com/pkg/errors"
)

// ListStyle is a stripped down metadata version of Style returned when listing styles.
// https://docs.mapbox.com/api/maps/#list-styles
type ListStyle struct {
	Version  int64     `json:"version,omitempty"`
	Name     string    `json:"name,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	ID       string    `json:"id,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	Owner    string    `json:"owner,omitempty"`
}

// https://docs.mapbox.com/mapbox-gl-js/style-spec/root/
// https://docs.mapbox.com/api/maps/#the-style-object
type Style struct {
	// Style specification version number. Must be 8.
	Version int `json:"version,omitempty"`
	// A human-readable name for the style.
	Name string `json:"name,omitempty"`
	// Arbitrary properties useful to track with the stylesheet, but do not influence rendering.
	// Properties should be prefixed to avoid collisions, like 'mapbox:'.
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	// Default map center in longitude and latitude.
	// The style center will be used only if the map has not been positioned by other means
	// (e.g. map options or user interaction).
	Center *[2]float64 `json:"center,omitempty"`
	// Default zoom level.
	// The style zoom will be used only if the map has not been positioned by other means
	// (e.g. map options or user interaction).
	Zoom *float64 `json:"zoom,omitempty"`
	// Default bearing, in degrees.
	// The bearing is the compass direction that is "up";
	// for example, a bearing of 90° orients the map so that east is up.
	// This value will be used only if the map has not been positioned by other means
	// (e.g. map options or user interaction).
	Bearing *float64 `json:"bearing,omitempty"`
	// Default pitch, in degrees.
	// Zero is perpendicular to the surface, for a look straight down at the map,
	// while a greater value like 60 looks ahead towards the horizon.
	// The style pitch will be used only if the map has not been positioned by other means
	// (e.g. map options or user interaction).
	Pitch *int64 `json:"pitch,omitempty"`
	// The global light source.
	Light *Light `json:"light,omitempty"`
	// Data source specifications.
	Sources map[string]json.RawMessage `json:"sources,omitempty"`
	// A base URL for retrieving the sprite image and metadata.
	// The extensions .png, .json and scale factor @2x.png will be automatically appended.
	// This property is required if any layer uses the
	// background-pattern, fill-pattern, line-pattern, fill-extrusion-pattern, or icon-image properties.
	// The URL must be absolute, containing the scheme, authority and path components.
	Sprite *string `json:"sprite,omitempty"`
	// A URL template for loading signed-distance-field glyph sets in PBF format.
	// The URL must include {fontstack} and {range} tokens.
	// This property is required if any layer uses the text-field layout property.
	// The URL must be absolute, containing the scheme, authority and path components.
	// For example "mapbox://fonts/mapbox/{fontstack}/{range}.pbf"
	Glyphs *string `json:"glyphs,omitempty"`
	// Layers will be drawn in the order of this array.
	Layers []json.RawMessage `json:"layers,omitempty"`
	// A global transition definition to use as a default across properties,
	// to be used for timing transitions between one value and the next when no property-specific transition is set.
	// Collision-based symbol fading is controlled independently of the style's transition property.
	Transition *Transition `json:"transition,omitempty"`
	// The date and time the style was created.
	Created time.Time `json:"created,omitempty"`
	// The ID of the style.
	Id string `json:"id,omitempty"`
	// The date and time the style was last modified.
	Modified time.Time `json:"modified,omitempty"`
	// The username of the style owner.
	Owner string `json:"owner,omitempty"`
	// Access control for the style, either public or private.
	// Private styles require an access token belonging to the owner.
	// Public styles may be requested with an access token belonging to any user.
	Visibility string
}

func (c *Client) ListStyles(draft bool) ([]ListStyle, error) {

	url := c.baseURL
	url.Path = path.Join(url.Path, "styles/v1/", c.username)
	q := url.Query()
	if draft {
		q.Add("draft", "true")
	} else {
		q.Add("draft", "false")
	}
	url.RawQuery = q.Encode()

	var allStyles []ListStyle

	requestURL := &url

	for requestURL != nil {

		var styles []ListStyle

		resp, err := c.do("GET", *requestURL, nil, &styles)
		if err != nil {
			return nil, errors.Wrap(err, "making request")
		}
		requestURL = c.nextPageURL(resp.Header)

		allStyles = append(allStyles, styles...)
	}

	return allStyles, nil
}

func (c *Client) GetStyle(id string, draft bool) (Style, error) {

	url := c.baseURL
	url.Path = path.Join(url.Path, "styles/v1/", c.username, id)
	if draft {
		url.Path = path.Join(url.Path, "draft")
	}

	var style Style

	_, err := c.do("GET", url, nil, &style)
	if err != nil {
		return Style{}, errors.Wrap(err, "making request")
	}

	return style, nil
}
