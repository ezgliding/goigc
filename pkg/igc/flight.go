// Copyright The ezgliding Authors.
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
	CompetitionID  string
	CompetitionURL string
	Speed          float64
	Comments       string
}
