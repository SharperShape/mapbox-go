package mapbox

import (
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// TilesetType is the data type of the tileset, either "vector" or "raster".
type TilesetType string

const (
	VectorTileset TilesetType = "vector"
	RasterTileset TilesetType = "raster"
)

// TilesetVisibility is either "public" or "private".
type TilesetVisibility string

const (
	PublicTileset  TilesetVisibility = "public"
	PrivateTileset TilesetVisibility = "private"
)

// ListTilesetsParams defines optional parameters for the ListTilesets request.
type ListTilesetsParams struct {
	Type       *TilesetType
	Visibility *TilesetVisibility
	SortBy     *SortBy
	Limit      *int
}

type Tileset struct {
	Type        string     `json:"type,omitempty"`
	Center      [3]float64 `json:"center,omitempty"`
	Created     time.Time  `json:"created,omitempty"`
	Description string     `json:"description,omitempty"`
	Filesize    int64      `json:"filesize,omitempty"`
	ID          string     `json:"id,omitempty"`
	Modified    time.Time  `json:"modified,omitempty"`
	Name        string     `json:"name,omitempty"`
	Visibility  string     `json:"visibility,omitempty"`
	Status      string     `json:"status,omitempty"`
}

func (c *Client) ListTilesets(params ListTilesetsParams) ([]Tileset, error) {

	url := c.baseURL
	url.Path = path.Join(url.Path, "tilesets/v1/", c.username)
	q := url.Query()
	if params.Type != nil {
		q.Add("type", string(*params.Type))
	}
	if params.Visibility != nil {
		q.Add("visibility", string(*params.Visibility))
	}
	if params.SortBy != nil {
		q.Add("sortby", string(*params.SortBy))
	}
	if params.Limit != nil {
		q.Add("limit", strconv.Itoa(*params.Limit))
	}
	url.RawQuery = q.Encode()

	var allTilesets []Tileset

	requestURL := &url

	for requestURL != nil {

		var tilesets []Tileset

		resp, err := c.do("GET", *requestURL, nil, &tilesets)
		if err != nil {
			return nil, errors.Wrap(err, "making request")
		}
		requestURL = c.nextPageURL(resp.Header)

		allTilesets = append(allTilesets, tilesets...)
	}

	return allTilesets, nil
}
