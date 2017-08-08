// Copyright 2014 The ezgliding Authors.
//
// This file is part of ezgliding.
//
// ezgliding is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// ezgliding is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ezgliding.  If not, see <http://www.gnu.org/licenses/>.
//
// Author: Ricardo Rocha <rocha.porto@gmail.com>

package igc

import (
	"strconv"
	"github.com/golang/geo/s2"
	"time"
)

const (
	// EarthRadius ...
	EarthRadius = 6371.0
)

// Point represents a gps read (single point in the track).
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

// NewPoint creates a new Point struct and returns it.
// It initializes all structures to zero values.
func NewPoint() Point {
	return NewPointFromLatLng(0, 0)
}

// NewPointFromLatLng ...
func NewPointFromLatLng(lat float64, lng float64) Point {
	return Point {
		LatLng: s2.LatLngFromDegrees(lat, lng),
		IData: make(map[string]string),
	}
}

// NewPointFromDMS returns a Point corresponding to the given string.
func NewPointFromDMS(lat string, lng string) Point {
	return NewPointFromLatLng(
		DecimalFromDMS(lat), DecimalFromDMS(lng),
	)
}

// DecimalFromDMS returns the decimal value corresponding to the given coordinates.
// The coordinates are expected in Degrees, Minutes, Seconds format.
func DecimalFromDMS(dms string) float64 {
	var degrees, minutes, seconds float64
	if len(dms) == 7 {
		degrees, _ = strconv.ParseFloat(dms[1:3], 64)
		minutes, _ = strconv.ParseFloat(dms[3:5], 64)
		seconds, _ = strconv.ParseFloat(dms[5:], 64)
	} else {
		degrees, _ = strconv.ParseFloat(dms[1:4], 64)
		minutes, _ = strconv.ParseFloat(dms[4:6], 64)
		seconds, _ = strconv.ParseFloat(dms[6:], 64)
	}
	var r float64
	r = degrees + (minutes / 60.0) + (seconds / 3600.0)
	if dms[0] == 'S' || dms[0] == 'W' {
		r = r * -1
	}
	return r
}

// NewPointFromDMD returns a Point corresponding to the given string.
func NewPointFromDMD(lat string, lng string) Point {
	return NewPointFromLatLng(
		DecimalFromDMD(lat), DecimalFromDMD(lng),
	)
}

// DecimalFromDMD ...
func DecimalFromDMD(dmd string) float64 {
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

// Distance dd
func (p *Point) Distance(b Point) float64 {
	return float64(p.Distance(b) * EarthRadius)
}
