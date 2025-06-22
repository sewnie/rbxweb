// package legacy selectively implements undocumented legacy Roblox Web APIs
package legacy

import (
	"github.com/apprehensions/rbxweb/internal/api"
	"github.com/apprehensions/rbxweb/services/games"
	"github.com/apprehensions/rbxweb/services/universes"
)

// LegacyServiceV1 handles legacy Roblox Web APIs.
type LegacyServiceV1 api.Service

// GetPlaceUniverse returns the given Place's Universe ID.
// If an error occurs, the returned Universe ID will be 0.
//
// Uses the undocumented API service universes/v1.
func (l *LegacyServiceV1) GetPlaceUniverse(pid games.PlaceID) (universes.UniverseID, error) {
	r := struct {
		UniverseID universes.UniverseID `json:"universeId"`
	}{}

	path := l.Client.Path("universes/v1/places/%d/universe", nil, pid)
	err := l.Client.Execute("GET", "apis", path, nil, &r)
	if err != nil {
		return 0, err
	}

	return r.UniverseID, nil
}
