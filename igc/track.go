// Copyright ©2015 Ricardo Rocha <rocha.porto@gmail.com>
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

	"github.com/kellydunn/golang-geo"
)

// Track holds all IGC flight data (header and gps track).
type Track struct {
	Header
	Points        []Point
	K             []K
	Events        []Event
	Satellites    []Satellite
	Logbook       []LogEntry
	Task          Task
	DGPSStationID string
	Signature     string
}

// NewTrack returns a new instance of Track.
// It initializes all the structures with zero values.
func NewTrack() Track {
	track := Track{}
	return track
}

// Header holds the meta information of a track.
type Header struct {
	Manufacturer     string
	UniqueID         string
	AdditionalData   string
	Date             time.Time
	FixAccuracy      int64
	Pilot            string
	Crew             string
	GliderType       string
	GliderID         string
	GPSDatum         string
	FirmwareVersion  string
	HardwareVersion  string
	FlightRecorder   string
	GPS              string
	PressureSensor   string
	CompetitionID    string
	CompetitionClass string
	Timezone         time.Location
}

// Point represents a gps read (single point in the track).
type Point struct {
	geo.Point
	Time             time.Time
	FixValidity      byte
	PressureAltitude int64
	GNSSAltitude     int64
	IData            map[string]string
	NumSatellites    int
	Description      string
}

// NewPoint creates a new Point struct and returns it.
// It initializes all structures to zero values.
func NewPoint() Point {
	var pt Point
	pt.IData = make(map[string]string)
	return pt
}

type K struct {
	Time   time.Time
	Fields map[string]string
}

type Satellite struct {
	Time time.Time
	Ids  []string
}

type Event struct {
	Time time.Time
	Type string
	Data string
}

// Task is a pre-declared task to be performed.
type Task struct {
	DeclarationDate time.Time
	Date            time.Time
	Number          int
	Takeoff         Point
	Start           Point
	Turnpoints      []Point
	Finish          Point
	Landing         Point
	Description     string
}

// LogEntry holds a logbook/comment entry in the IGC file.
type LogEntry struct {
	Type string
	Text string
}

// Manufacturer holds the char identifier, the short id and the full name of
// an IGC Manufacturer, as defined in Appendix A (Codes for Manufacturers)
// of the IGC spec.
type Manufacturer struct {
	char  byte
	short string
	name  string
}

// Manufacturers holds the list of available manufacturers, as defined in
// Appendix A (Codes for Manufacturers) of the IGC spec.
var Manufacturers = map[string]Manufacturer{
	"GCS": Manufacturer{'A', "GCS", "Garrecht"},
	"LGS": Manufacturer{'B', "LGS", "Logstream"},
	"CAM": Manufacturer{'C', "CAM", "Cambridge Aero Instruments"},
	"DSX": Manufacturer{'D', "DSX", "Data Swan/DSX"},
	"EWA": Manufacturer{'E', "EWA", "EW Avionics"},
	"FIL": Manufacturer{'F', "FIL", "Filser"},
	"FLA": Manufacturer{'G', "FLA", "Flarm (Track Alarm)"},
	"SCH": Manufacturer{'H', "SCH", "Scheffel"},
	"ACT": Manufacturer{'I', "ACT", "Aircotec"},
	"CNI": Manufacturer{'K', "CNI", "ClearNav Instruments"},
	"NKL": Manufacturer{'K', "NKL", "NKL"},
	"LXN": Manufacturer{'L', "LXN", "LX Navigation"},
	"IMI": Manufacturer{'M', "IMI", "IMI Gliding Equipment"},
	"NTE": Manufacturer{'N', "NTE", "New Technologies s.r.l."},
	"NAV": Manufacturer{'O', "NAV", "Naviter"},
	"PES": Manufacturer{'P', "PES", "Peschges"},
	"PRT": Manufacturer{'R', "PRT", "Print Technik"},
	"SDI": Manufacturer{'S', "SDI", "Streamline Data Instruments"},
	"TRI": Manufacturer{'T', "TRI", "Triadis Engineering GmbH"},
	"LXV": Manufacturer{'V', "LXV", "LXNAV d.o.o."},
	"WES": Manufacturer{'W', "WES", "Westerboer"},
	"XYY": Manufacturer{'X', "XYY", "Other manufacturer"},
	"ZAN": Manufacturer{'Z', "ZAN", "Zander"},
}
