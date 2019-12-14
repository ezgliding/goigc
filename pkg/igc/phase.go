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
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/golang/geo/s2"
	kml "github.com/twpayne/go-kml"
)

// PhaseType represents a flight phase.
//
// Possible values include Towing, PossibleCruising/Cruising,
// PossibleCircling/Circling, Unknown.
type PhaseType int

const (
	Unknown          PhaseType = 0
	Towing           PhaseType = 1
	PossibleCruising PhaseType = 2
	Cruising         PhaseType = 3
	PossibleCircling PhaseType = 4
	Circling         PhaseType = 5
)

// CirclingType indicates Left, Right or Mixed circling.
type CirclingType int

const (
	Mixed CirclingType = 0
	Left  CirclingType = 1
	Right CirclingType = 2
)

const (
	// MinTurnRate is the min rate to consider circling
	MinTurnRate = 6.5
	// MaxTurnRate is the max rate considered for valid turns
	MaxTurnRate = 22.5
	// MinCirclingTime is used to decide when a switch to circling occurs.
	// This value is used when calculating flight phases to switch from
	// PossibleCircling to Circling.
	MinCirclingTime = 15
	// MinCruisingTime is used to decide when a switch to cruising occurs.
	// This value is used when calculating flight phases to switch from
	// PossibleCruising to Cruising.
	MinCruisingTime = 10
)

// Phase is a flight phase (towing, cruising, circling).
type Phase struct {
	Type         PhaseType
	CirclingType CirclingType
	Start        Point
	StartIndex   int
	End          Point
	EndIndex     int
	AvgVario     float64
	TopVario     float64
	AvgGndSpeed  float64
	TopGndSpeed  float64
	Distance     float64
	LD           float64
	Centroid     s2.LatLng
	CellID       s2.CellID
}

// Phases returns the list of flight phases for the Track.
// Each phases is one of Cruising, Circling, Towing or Unknown.
func (track *Track) Phases() ([]Phase, error) {

	if len(track.phases) > 0 {
		return track.phases, nil
	}

	var currPhase PhaseType
	var startIndex int
	var currPoint Point
	var turning bool
	//var turnRate float64

	currPhase = Cruising
	track.phases = []Phase{
		Phase{Type: Cruising, StartIndex: 0, Start: track.Points[0]},
	}

	// we need the bearings for each point to calculate turn rates
	var d float64
	for i := 1; i < len(track.Points); i++ {
		track.Points[i-1].bearing = track.Points[i-1].Bearing(track.Points[i])
		d = track.Points[i-1].Distance(track.Points[i])
		track.Points[i].distance = track.Points[i-1].distance + d
		track.Points[i].speed = d / track.Points[i].Time.Sub(track.Points[i-1].Time).Seconds()
	}

	for i := 0; i < len(track.Points)-1; i++ {
		currPoint = track.Points[i]
		turning, _ = track.isTurning(i)

		if currPhase == Cruising {
			// if cruising check for turning
			if turning {
				// set possible circling if turning
				currPhase = PossibleCircling
				startIndex = i
			} // else continue
		} else if currPhase == PossibleCircling {
			// if possible circling check for turning longer than min circling time
			if turning {
				if currPoint.Time.Sub(track.Points[startIndex].Time).Seconds() > MinCirclingTime {
					// if true then set circling
					currPhase = Circling
					track.wrapPhase(startIndex, Circling)
				}
			} else {
				// if not go back to cruising
				currPhase = Cruising
			}
		} else if currPhase == Circling {
			// if circling check for stopping to turn
			if !turning {
				// if stopping set possible cruising
				currPhase = PossibleCruising
				startIndex = i
			}
		} else if currPhase == PossibleCruising {
			// if possible cruising check for longer than min cruising
			if !turning {
				if currPoint.Time.Sub(track.Points[startIndex].Time).Seconds() > MinCruisingTime {
					// if true then set cruising
					currPhase = Cruising
					track.wrapPhase(startIndex, Cruising)
				}
			} else {
				// if not go back to circling
				currPhase = Circling
			}
		}
	}

	return track.phases, nil
}

func (track *Track) wrapPhase(index int, phaseType PhaseType) {
	p := &track.phases[len(track.phases)-1]

	p.EndIndex = index
	p.End = track.Points[index]

	// compute phase stats
	altGain := float64(p.End.GNSSAltitude - p.Start.GNSSAltitude)
	p.Distance = p.End.distance - p.Start.distance
	p.AvgVario = altGain / p.Duration().Seconds()
	p.AvgGndSpeed = p.Distance / (p.Duration().Seconds() / 3600)

	if p.Type == Cruising {
		p.LD = p.Distance * 1000.0 / math.Abs(altGain)
	}
	pts := make([]s2.LatLng, p.EndIndex-p.StartIndex)
	for i := p.StartIndex; i < p.EndIndex; i++ {
		pts[i-p.StartIndex] = track.Points[i].LatLng
	}
	centroid := s2.LatLngFromPoint(s2.PolylineFromLatLngs(pts).Centroid())
	p.CellID = s2.CellIDFromLatLng(centroid)
	p.Centroid = centroid

	track.phases = append(track.phases, Phase{Type: phaseType, StartIndex: index, Start: track.Points[index]})
}

func (track *Track) isTurning(i int) (bool, float64) {
	turnRate := (track.Points[i+1].bearing - track.Points[i].bearing).Abs().Degrees() / track.Points[i+1].Time.Sub(track.Points[i].Time).Seconds()
	return math.Abs(turnRate) > MinTurnRate, turnRate
}

// Duration returns the duration of this flight phase.
func (p *Phase) Duration() time.Duration {
	return p.End.Time.Sub(p.Start.Time)
}

func (track *Track) encodePhasesKML() (*kml.CompoundElement, error) {

	result := kml.Document()
	result.Add(
		kml.SharedStyle(
			"cruising",
			kml.LineStyle(
				kml.Color(color.RGBA{R: 0, G: 0, B: 255, A: 127}),
				kml.Width(4),
			),
		),
		kml.SharedStyle(
			"circling",
			kml.LineStyle(
				kml.Color(color.RGBA{R: 0, G: 255, B: 0, A: 127}),
				kml.Width(4),
			),
		),
		kml.SharedStyle(
			"attempt",
			kml.LineStyle(
				kml.Color(color.RGBA{R: 255, G: 0, B: 0, A: 127}),
				kml.Width(4),
			),
		),
	)

	phases, err := track.Phases()
	if err != nil {
		return result, err
	}

	for i := 0; i < len(phases)-2; i++ {
		phase := phases[i]
		//fmt.Printf("%v\t%v\n", ph.Start, ph.End)
		coords := make([]kml.Coordinate, phase.EndIndex-phase.StartIndex+1)
		for i := phase.StartIndex; i <= phase.EndIndex; i++ {
			p := track.Points[i]
			coords[i-phase.StartIndex].Lat = p.Lat.Degrees()
			coords[i-phase.StartIndex].Lon = p.Lng.Degrees()
			coords[i-phase.StartIndex].Alt = float64(p.GNSSAltitude)
		}
		style := "#cruising"
		if phase.Type == Circling && phase.End.Time.Sub(phase.Start.Time).Seconds() < 45 {
			style = "#attempt"
		} else if phase.Type == Circling {
			style = "#circling"
		}
		result.Add(
			kml.Placemark(
				kml.StyleURL(style),
				kml.LineString(
					kml.Extrude(false),
					kml.Tessellate(false),
					kml.AltitudeMode("absolute"),
					kml.Coordinates(coords...),
				),
			))

		name := fmt.Sprintf("Lat: %v Lng: %v",
			phase.Centroid.Lat.Degrees(), phase.Centroid.Lng.Degrees())
		desc := fmt.Sprintf("Alt Gain: %dm (%dm %dm)<br/>Distance: %.2fkm<br/>Speed: %.2fkm/h<br/>LD: %v<br/>Vario: %.1fm/s<br/>Cell: %v<br/>",
			phase.End.GNSSAltitude-phase.Start.GNSSAltitude,
			phase.Start.GNSSAltitude, phase.End.GNSSAltitude, phase.Distance,
			phase.AvgGndSpeed, phase.LD, phase.AvgVario, phase.CellID)
		result.Add(
			kml.Placemark(
				kml.Name(name),
				kml.Description(desc),
				kml.Point(
					kml.Coordinates(kml.Coordinate{
						Lon: phase.Centroid.Lng.Degrees(), Lat: phase.Centroid.Lat.Degrees(),
					}),
				),
			))
	}
	return result, nil
}
