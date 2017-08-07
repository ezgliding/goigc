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

package spatial

import (
	"reflect"
	"testing"
)

type DMS2DecimalTest struct {
	t  string
	in string
	r  float64
}

var dms2DecimalTests = []DMS2DecimalTest{
	{
		"latitude north conversion",
		"N323200",
		32.53333333333333,
	},
	{
		"latitude south conversion",
		"S323200",
		-32.53333333333333,
	},
	{
		"longitude east conversion",
		"E1002233",
		100.37583333333333,
	},
	{
		"longitude west conversion",
		"W1002233",
		-100.37583333333333,
	},
}

func TestDMS2Decimal(t *testing.T) {
	for _, test := range dms2DecimalTests {
		result := DMS2Decimal(test.in)
		if result != test.r {
			t.Errorf("test %v failed, expected %v got %v", test.t, test.r, result)
			continue
		}
	}
}

type DMD2DecimalTest struct {
	t  string
	in string
	r  float64
}

var dmd2DecimalTests = []DMD2DecimalTest{
	{
		"latitude north conversion",
		"N4616018",
		46.26696666666667,
	},
	{
		"latitude north conversion inverted",
		"4616018N",
		46.26696666666667,
	},
	{
		"latitude south conversion",
		"S4616018",
		-46.26696666666667,
	},
	{
		"latitude south conversion inverted",
		"4616018S",
		-46.26696666666667,
	},
	{
		"longitude east conversion",
		"E00627679",
		6.461316666666667,
	},
	{
		"longitude east conversion inverted",
		"00627679E",
		6.461316666666667,
	},
	{
		"longitude west conversion",
		"W00627679",
		-6.461316666666667,
	},
	{
		"longitude west conversion inverted",
		"00627679W",
		-6.461316666666667,
	},
}

func TestDMD2Decimal(t *testing.T) {
	for _, test := range dmd2DecimalTests {
		result := DMD2Decimal(test.in)
		if result != test.r {
			t.Errorf("test %v failed, expected %v got %v", test.t, test.r, result)
			continue
		}
	}
}
