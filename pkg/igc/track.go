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
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"time"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	kml "github.com/twpayne/go-kml"
	"gopkg.in/yaml.v3"
)

// MaxSpeed is the maximum theoretical speed for a glider.
//
// It is used to detect bad GPS coordinates, which should be removed from the track.
const MaxSpeed float64 = 500.0

// Track holds all IGC flight data (header and GPS points).
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

// NewTrack returns a new instance of Track, with fields initialized to zero.
func NewTrack() Track {
	track := Track{}
	return track
}

// Header holds the meta information of a track.
//
// This is the H record in the IGC specification, section A3.2.
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
	Timezone         int
}

// K holds flight data needed less often than Points.
//
// This is the K record in the IGC specification, section A4.4. Fields
// is a map between a given content type and its value, with the possible
// content types being defined in the J record.
//
// Examples of content types include heading true (HDT) or magnetic (HDM),
// airspeed (IAS), etc.
type K struct {
	Time   time.Time
	Fields map[string]string
}

// Satellite holds the IDs of the available satellites at a given Time.
//
// This is the F record in the IGC specification, section A4.3.
type Satellite struct {
	Time time.Time
	Ids  []string
}

// Event holds data records triggered at a given time.
//
// This is the E record in the IGC specification, section A4.2. The events
// can be pilot initiated (with a PEV code), proximity alerts, etc.
type Event struct {
	Time time.Time
	Type string
	Data string
}

// Task holds all the metadata put in a pre-declared task to be performed.
//
// This is the C record in the IGC specification, section A3.5.
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

// Distance returns the total distance in kms between the turn points.
//
// It includes the Start and Finish fields as the first and last point,
// respectively, with the Turnpoints in the middle. The return value is
// sum of all distances between each consecutive point.
func (task *Task) Distance() float64 {
	d := 0.0
	p := []Point{task.Start}
	p = append(p, task.Turnpoints...)
	p = append(p, task.Finish)
	for i := 0; i < len(p)-1; i++ {
		d += p[i].Distance(p[i+1])
	}
	return d
}

// LogEntry holds a logbook/comment entry, in free format.
//
// This is the L record in the IGC specification, section A4.5.
type LogEntry struct {
	Type string
	Text string
}

// Manufacturer holds manufacturer name, short ID and char identifier.
//
// The list of manufacturers is defined in the IGC specification,
// section A2.5.6. A map Manufacturers is available in this library.
type Manufacturer struct {
	char  byte   //nolint
	short string //nolint
	name  string //nolint
}

// Manufacturers holds the list of available manufacturers.
//
// This list is defined in the IGC specification, section A2.5.6.
var Manufacturers = map[string]Manufacturer{
	"GCS": {'A', "GCS", "Garrecht"},
	"LGS": {'B', "LGS", "Logstream"},
	"CAM": {'C', "CAM", "Cambridge Aero Instruments"},
	"DSX": {'D', "DSX", "Data Swan/DSX"},
	"EWA": {'E', "EWA", "EW Avionics"},
	"FIL": {'F', "FIL", "Filser"},
	"FLA": {'G', "FLA", "Flarm (Track Alarm)"},
	"SCH": {'H', "SCH", "Scheffel"},
	"ACT": {'I', "ACT", "Aircotec"},
	"CNI": {'K', "CNI", "ClearNav Instruments"},
	"NKL": {'K', "NKL", "NKL"},
	"LXN": {'L', "LXN", "LX Navigation"},
	"IMI": {'M', "IMI", "IMI Gliding Equipment"},
	"NTE": {'N', "NTE", "New Technologies s.r.l."},
	"NAV": {'O', "NAV", "Naviter"},
	"PES": {'P', "PES", "Peschges"},
	"PRT": {'R', "PRT", "Print Technik"},
	"SDI": {'S', "SDI", "Streamline Data Instruments"},
	"TRI": {'T', "TRI", "Triadis Engineering GmbH"},
	"LXV": {'V', "LXV", "LXNAV d.o.o."},
	"WES": {'W', "WES", "Westerboer"},
	"XYY": {'X', "XYY", "Other manufacturer"},
	"ZAN": {'Z', "ZAN", "Zander"},
}

func (track *Track) Cleanup() (Track, error) {
	clean := *track

	i := 1
	for i < len(clean.Points) {
		if clean.Points[i-1].Speed(clean.Points[i]) > MaxSpeed {
			clean.Points = append(clean.Points[:i], clean.Points[i+1:]...)
		}
		i = i + 1
	}
	return clean, nil
}

func (track *Track) Simplify(tolerance float64) (Track, error) {
	r := polylineFromPoints(track.Points).SubsampleVertices(s1.Angle(tolerance))
	points := make([]Point, len(r))
	for i, v := range r {
		points[i] = track.Points[v]
	}

	simplified := *track
	simplified.Points = points
	return simplified, nil
}

func (track *Track) Encode(format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(track, "", "  ")
	case "kml", "kmz":
		return encodeKML(track, format)
	case "yaml":
		return yaml.Marshal(track)
	default:
		return []byte{}, fmt.Errorf("unsupported format '%v'", format)
	}
}

func encodeKML(track *Track, format string) ([]byte, error) {

	coords := make([]kml.Coordinate, len(track.Points))
	for i, p := range track.Points {
		coords[i].Lat = p.Lat.Degrees()
		coords[i].Lon = p.Lng.Degrees()
		coords[i].Alt = float64(p.GNSSAltitude)
	}

	metadata := fmt.Sprintf("%v : %v : %v", track.Date, track.Pilot, track.GliderType)

	k := kml.KML(
		kml.Document(
			kml.Name(metadata),
			kml.Description(""),
			kml.SharedStyle(
				"yellowLineGreenPoly",
				kml.LineStyle(
					kml.Color(color.RGBA{R: 255, G: 255, B: 0, A: 127}),
					kml.Width(4),
				),
			),
			kml.Placemark(
				kml.Name(metadata),
				kml.StyleURL("#yellowLineGreenPoly"),
				kml.LineString(
					kml.Extrude(false),
					kml.Tessellate(false),
					kml.AltitudeMode("absolute"),
					kml.Coordinates(coords...),
				),
			),
		),
	)
	buf := new(bytes.Buffer)
	if err := k.WriteIndent(buf, "", "  "); err != nil {
		return buf.Bytes(), err
	}

	if format == "kmz" {
		zipbuf := new(bytes.Buffer)

		w := zip.NewWriter(zipbuf)
		f, err := w.Create("flight.kml")
		if err != nil {
			return []byte{}, err
		}
		_, err = f.Write(buf.Bytes())
		if err != nil {
			return []byte{}, err
		}
		err = w.Close()
		if err != nil {
			return []byte{}, err
		}

		return zipbuf.Bytes(), nil
	}
	return buf.Bytes(), nil
}

func polylineFromPoints(points []Point) *s2.Polyline {
	p := make(s2.Polyline, len(points))
	for k, v := range points {
		p[k] = s2.PointFromLatLng(v.LatLng)
	}
	return &p
}
