package main

import (
	"encoding/json"
	"io"
)

func NewLeague(r io.Reader) ([]Player, error) {
	var league []Player
	if err := json.NewDecoder(r).Decode(&league); err != nil {
		return nil, err
	}
	return league, nil
}
