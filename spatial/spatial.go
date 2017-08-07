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

// Package spatial provides functionality for handling spatial data.
//
// This includes conversion for lat/lon between different formats (dms,
// decimal, ...) and other functions coming in the future.
package spatial

import (
	"strconv"
)

// DMS2Decimal converts the given coordinates from DMS to decimal format.
func DMS2Decimal(dms string) float64 {
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

// DMD2Decimal converts the given coordinates from DMD (deg,min,decimalmin) to decimal format.
func DMD2Decimal(dmd string) float64 {
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
