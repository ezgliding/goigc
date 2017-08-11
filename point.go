// Copyright @2017 The ezgliding Authors.
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
	"strconv"
	"time"

	"github.com/golang/geo/s2"
)

const (
	// EarthRadius is the average earth radius.
	EarthRadius = 6371.0
)

// Point represents a GPS recording (single point in the track).
//
// It is based on a golang-geo s2 LatLng, adding extra metadata such as
// the Time the point was recorded, pressure and GNSS altitude, number of
// satellites available and extra metadata added by the recorder.
//
// You can use all methods available for a s2.LatLng on this struct.
type Point struct {
	s2.LatLng
	Time             time.Time
	FixValidity      byte
	PressureAltitude int64
	GNSSAltitude     int64
	IData            map[string]string
	NumSatellites    int
	Description      string
}

// NewPoint returns a new Point set to latitude and longitude 0.
func NewPoint() Point {
	return NewPointFromLatLng(0, 0)
}

// NewPointFromLatLng returns a new Point with the given latitude and longitude.
func NewPointFromLatLng(lat float64, lng float64) Point {
	return Point{
		LatLng: s2.LatLngFromDegrees(lat, lng),
		IData:  make(map[string]string),
	}
}

// NewPointFromDMS returns a Point corresponding to the given string in DMS format.
//
// DecimalFromDMS includes more information regarding this format.
func NewPointFromDMS(lat string, lng string) Point {
	return NewPointFromLatLng(
		DecimalFromDMS(lat), DecimalFromDMS(lng),
	)
}

// DecimalFromDMS returns the decimal representation (in radians) of the given DMS.
//
// DMS is a representation of a coordinate in Decimal,Minutes,Seconds, with an
// extra character indicating north, south, east, west.
//
// Examples: N512646, W0064312, S342244, E0021233
func DecimalFromDMS(dms string) float64 {
	var degrees, minutes, seconds float64
	if len(dms) == 7 {
		degrees, _ = strconv.ParseFloat(dms[1:3], 64)
		minutes, _ = strconv.ParseFloat(dms[3:5], 64)
		seconds, _ = strconv.ParseFloat(dms[5:], 64)
	} else if len(dms) == 8 {
		degrees, _ = strconv.ParseFloat(dms[1:4], 64)
		minutes, _ = strconv.ParseFloat(dms[4:6], 64)
		seconds, _ = strconv.ParseFloat(dms[6:], 64)
	} else {
		return 0
	}
	var r float64
	r = degrees + (minutes / 60.0) + (seconds / 3600.0)
	if dms[0] == 'S' || dms[0] == 'W' {
		r = r * -1
	}
	return r
}

// NewPointFromDMD returns a Point corresponding to the given string in DMD format.
//
// DecimalFromDMD includes more information regarding this format.
func NewPointFromDMD(lat string, lng string) Point {
	return NewPointFromLatLng(
		DecimalFromDMD(lat), DecimalFromDMD(lng),
	)
}

// DecimalFromDMD returns the decimal representation (in radians) of the given DMD.
//
// DMD is a representation of a coordinate in Decimal,Minutes,100thMinute with an
// extra character indicating north, south, east, west.
//
// Examples: N512688, W0064364, S342212, E0021275
func DecimalFromDMD(dmd string) float64 {
	if len(dmd) != 8 && len(dmd) != 9 {
		return 0
	}

	var degrees, minutes, dminutes float64
	if dmd[0] == 'S' || dmd[0] == 'N' {
		degrees, _ = strconv.ParseFloat(dmd[1:3], 64)
		minutes, _ = strconv.ParseFloat(dmd[3:5], 64)
		dminutes, _ = strconv.ParseFloat(dmd[5:], 64)
	} else if dmd[len(dmd)-1] == 'S' || dmd[len(dmd)-1] == 'N' {
		degrees, _ = strconv.ParseFloat(dmd[0:2], 64)
		minutes, _ = strconv.ParseFloat(dmd[2:4], 64)
		dminutes, _ = strconv.ParseFloat(dmd[4:7], 64)
	} else if dmd[0] == 'W' || dmd[0] == 'E' {
		degrees, _ = strconv.ParseFloat(dmd[1:4], 64)
		minutes, _ = strconv.ParseFloat(dmd[4:6], 64)
		dminutes, _ = strconv.ParseFloat(dmd[6:], 64)
	} else if dmd[len(dmd)-1] == 'W' || dmd[len(dmd)-1] == 'E' {
		degrees, _ = strconv.ParseFloat(dmd[0:3], 64)
		minutes, _ = strconv.ParseFloat(dmd[3:5], 64)
		dminutes, _ = strconv.ParseFloat(dmd[5:8], 64)
	}
	var r float64
	r = degrees + ((minutes + (dminutes / 1000.0)) / 60.0)
	if dmd[0] == 'S' || dmd[0] == 'W' || dmd[len(dmd)-1] == 'S' || dmd[len(dmd)-1] == 'W' {
		r = r * -1
	}
	return r
}

// Distance returns the great circle distance in kms to the given point.
//
// Internally it uses the golang-geo s2 LatLng.Distance() method, but converts
// its result (an angle) to kms considering the constance EarthRadius.
func (p *Point) Distance(b Point) float64 {
	return float64(p.LatLng.Distance(b.LatLng) * EarthRadius)
}
