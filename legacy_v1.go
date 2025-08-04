package rbxweb

// LegacyServiceV1 partially handles legacy Roblox Web APIs.
type LegacyServiceV1 service

// GetPlaceUniverse returns the given Place's Universe ID.
// If an error occurs, the returned Universe ID will be 0.
//
// Uses the undocumented API service universes/v1.
func (l *LegacyServiceV1) GetPlaceUniverse(pid PlaceID) (UniverseID, error) {
	r := struct {
		UniverseID UniverseID `json:"universeId"`
	}{}

	path := Path("universes/v1/places/%d/universe", nil, pid)
	err := l.Client.Execute("GET", "apis", path, nil, &r)
	if err != nil {
		return 0, err
	}

	return r.UniverseID, nil
}
