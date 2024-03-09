package rbxweb

import (
	"strconv"
)

type (
	UniverseID int64
	CreatorID  int64
	PlaceID    int64
)

// GetPlaceUniverse returns the given Place's Universe ID.
//
// Uses the undocumented API service universes/v1
func GetPlaceUniverse(placeID PlaceID) (UniverseID, error) {
	r := struct {
		UniverseID `json:"universeId"`
	}{}

	pid := strconv.FormatInt(int64(placeID), 10)
	err := Request("GET",
		GetURL("apis", "universes/v1/places/"+pid+"/universe", nil), nil, &r)
	if err != nil {
		return 0, err
	}

	return r.UniverseID, nil
}
