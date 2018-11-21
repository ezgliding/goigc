// Copyright 2018 The ezgliding authors. All rights reserverd.

package igc

import (
	"time"
)

// Flight represents a flight submission in an online competition.
//
// It includes all the flight metadata available in the online competition,
// including computed data like speed or distance. It also includes a url to
// the flight track file (usually in IGC format).
type Flight struct {
	URL            string
	ID             string
	Pilot          string
	Club           string
	Date           time.Time
	Takeoff        string
	Region         string
	Country        string
	Distance       float64
	Points         float64
	Glider         string
	Type           string
	TrackURL       string
	TrackID        string
	CompetitionID  string `json:",omitempty"`
	CompetitionURL string `json:",omitempty"`
	Speed          float64
	Comments       string
}
