package application

import (
	"encoding/json"
	"io"
)

type League []Player

func NewLeague(r io.Reader) (League, error) {
	var league League
	if err := json.NewDecoder(r).Decode(&league); err != nil {
		return nil, err
	}
	return league, nil
}

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
