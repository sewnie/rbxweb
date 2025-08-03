// games selectively implements the 'games' Roblox Web API.
package games

import (
	"net/url"

	"github.com/sewnie/rbxweb/internal/api"
	"github.com/sewnie/rbxweb/internal/stringutils"
	"github.com/sewnie/rbxweb/services/universes"
)

// GamesServiceV1 handles the 'games/v1' Roblox Web API.
type GamesServiceV1 api.Service

// GetGamesDetail returns a list of the game details of each given Universe ID.
func (g *GamesServiceV1) ListGamesDetails(uids []universes.UniverseID) ([]GameDetail, error) {
	gdr := struct {
		Data []GameDetail `json:"data"`
	}{}

	query := url.Values{"universeIds": stringutils.FormatSlice(uids)}
	err := g.Client.Execute("GET", "games", api.Path("v1/games", query), nil, &gdr)
	if err != nil {
		return nil, err
	}

	return gdr.Data, nil
}

// GetGameDetail returns the given Universe ID's game details.
//
// If none are found, nil will be returned.
func (g *GamesServiceV1) GetGameDetail(uid universes.UniverseID) (*GameDetail, error) {
	return api.GetList(g.ListGamesDetails([]universes.UniverseID{uid}))
}

// ListPlacesDetail returns a list of the place details of each given Place ID.
func (g *GamesServiceV1) ListPlacesDetails(pids []PlaceID) ([]PlaceDetail, error) {
	var pds []PlaceDetail

	query := url.Values{"placeIds": stringutils.FormatSlice(pids)}
	err := g.Client.Execute("GET", "games", api.Path("v1/games/multiget-place-details", query), nil, &pds)
	if err != nil {
		return nil, err
	}

	return pds, nil
}

// GetPlaceDetails returns the given Place ID's Place details.
//
// If none are found, nil will be returned.
func (g *GamesServiceV1) GetPlaceDetail(placeID PlaceID) (*PlaceDetail, error) {
	return api.GetList(g.ListPlacesDetails([]PlaceID{placeID}))
}
