package mapbox

import (
	"path"
	"strconv"
	"strings"
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

type TileJSON struct {
	Bounds       [4]float64    `json:"bounds,omitempty"`
	Center       [3]float64    `json:"center,omitempty"`
	Created      int64         `json:"created,omitempty"`
	Format       string        `json:"format,omitempty"`
	MinZoom      int64         `json:"minzoom,omitempty"`
	MaxZoom      int64         `json:"maxzoom,omitempty"`
	Name         string        `json:"name,omitempty"`
	Scheme       string        `json:"scheme,omitempty"`
	TileJSON     string        `json:"tile_json,omitempty"`
	Tiles        []string      `json:"tiles,omitempty"`
	VectorLayers []VectorLayer `json:"vector_layers,omitempty"`
}

type VectorLayer struct {
	Description string            `json:"description,omitempty"`
	Fields      map[string]string `json:"fields,omitempty"`
	ID          string            `json:"id,omitempty"`
	MaxZoom     int64             `json:"maxzoom,omitempty"`
	MinZoom     int64             `json:"minzoom,omitempty"`
	Source      string            `json:"source,omitempty"`
	SourceName  string            `json:"source_name,omitempty"`
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

func (c *Client) GetTileJSON(tilesetIDs ...string) (TileJSON, error) {

	ids := strings.Join(tilesetIDs, ",")

	url := c.baseURL
	url.Path = path.Join(url.Path, "v4", ids) + ".json"

	var metadata TileJSON

	_, err := c.do("GET", url, nil, &metadata)
	if err != nil {
		return TileJSON{}, err
	}

	return metadata, nil
}
