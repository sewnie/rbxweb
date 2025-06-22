// thumbnails selectively implements the 'thumbnails' Roblox Web API.
package thumbnails

import (
	"net/url"
	"strconv"

	"github.com/apprehensions/rbxweb/internal/api"
	"github.com/apprehensions/rbxweb/internal/stringutils"
	"github.com/apprehensions/rbxweb/services/universes"
)

// ThumbnailsServiceV1 handles the 'thumbnails/v1' Roblox Web API.
type ThumbnailsServiceV1 api.Service

// ListGamesIcons returns a list of Thumbnails for the given list of universeIDs, based on the named policy,
// thumbnail size, thumbnail format, and whether the thumbnail is circular.
func (t *ThumbnailsServiceV1) ListGamesIcons(uids []universes.UniverseID, policy ReturnPolicy, size string, format ThumbnailFormat, circular bool) ([]Thumbnail, error) {
	r := struct {
		Data []Thumbnail `json:"data"`
	}{}

	path := api.Path("v1/games/icons", url.Values{
		"universeIds":  stringutils.FormatSlice(uids),
		"returnPolicy": {string(policy)},
		"size":         {size},
		"format":       {string(format)},
		"isCircular":   {strconv.FormatBool(circular)},
	})

	err := t.Client.Execute("GET", "thumbnails", path, nil, &r)
	if err != nil {
		return nil, err
	}

	return r.Data, nil
}

// GetGameIcons returns a Thumbnail for the given universeID, based on the named policy,
// thumbnail size, thumbnail format, and whether the thumbnail is circular.
//
// If none are found, nil will be returned.
func (t *ThumbnailsServiceV1) GetGameIcon(universeID universes.UniverseID, policy ReturnPolicy, size string, format ThumbnailFormat, circular bool) (*Thumbnail, error) {
	return api.GetList(t.ListGamesIcons([]universes.UniverseID{universeID}, policy, size, format, circular))
}
