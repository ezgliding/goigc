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
	"testing"
	"time"
)

type point struct {
	lat string
	lng string
	t   time.Time
}

type PointTest struct {
	p1 point
	p2 point
	d  float64
	b  float64
	s  float64
}

var dmdTests = []PointTest{
	{
		p1: point{lat: "", lng: ""},
		p2: point{lat: "", lng: ""},
		d:  0,
		b:  0,
	},
	{
		p1: point{lat: "5110179N", lng: "00102644W"},
		p2: point{lat: "5110179N", lng: "00102644W"},
		d:  0,
		b:  0,
	},
}

var tests = []PointTest{
	{
		p1: point{lat: "", lng: ""},
		p2: point{lat: "", lng: ""},
		d:  0,
		b:  0,
	},
	{
		p1: point{lat: "N500359", lng: "W0054253", t: getTime("10:00:00")},
		p2: point{lat: "N500359", lng: "W0054253", t: getTime("10:00:00")},
		d:  0,
		b:  0,
		s:  0,
	},
	{
		p1: point{lat: "N500359", lng: "W0054253", t: getTime("10:00:00")},
		p2: point{lat: "N583838", lng: "W0030412", t: getTime("10:30:00")},
		d:  968.8535467131387,
		b:  9.119818104504075,
		s:  968.8535467131387 / 0.5,
	},
	{
		p1: point{lat: "S270201", lng: "E0303722", t: getTime("12:00:00")},
		p2: point{lat: "N523838", lng: "W0030412", t: getTime("13:30:00")},
		d:  9443.596093743798,
		b:  -19.75024484768977,
		s:  9443.596093743798 / 1.5,
	},
}

func TestPointDistanceDMD(t *testing.T) {
	for _, test := range dmdTests {
		p1 := NewPointFromDMD(test.p1.lat, test.p1.lng)
		p2 := NewPointFromDMD(test.p2.lat, test.p2.lng)
		result := p1.Distance(p2)
		if result != test.d {
			t.Errorf("p1: %v p2: %v :: expected distance %v got %+v", test.p1, test.p2, test.d, result)
			continue
		}
	}
}

func TestPointSpeed(t *testing.T) {
	for _, test := range tests {
		p1 := NewPointFromDMS(test.p1.lat, test.p1.lng)
		p1.Time = test.p1.t
		p2 := NewPointFromDMS(test.p2.lat, test.p2.lng)
		p2.Time = test.p2.t
		result := p1.Speed(p2)
		if result != test.s {
			t.Errorf("p1: %v p2: %v :: expected speed %v got %+v", test.p1, test.p2, test.s, result)
			continue
		}
	}
}

func TestPointDistance(t *testing.T) {
	for _, test := range tests {
		p1 := NewPointFromDMS(test.p1.lat, test.p1.lng)
		p2 := NewPointFromDMS(test.p2.lat, test.p2.lng)
		result := p1.Distance(p2)
		if result != test.d {
			t.Errorf("p1: %v p2: %v :: expected distance %v got %+v", test.p1, test.p2, test.d, result)
			continue
		}
	}
}

func TestPointBearing(t *testing.T) {
	for _, test := range tests {
		p1 := NewPointFromDMS(test.p1.lat, test.p1.lng)
		p2 := NewPointFromDMS(test.p2.lat, test.p2.lng)
		result := p1.Bearing(p2).Degrees()
		if result != test.b {
			t.Errorf("p1: %v p2: %v :: expected bearing %v got %+v", test.p1, test.p2, test.b, result)
			continue
		}
	}
}

func getTime(v string) time.Time {
	t, _ := time.Parse("15:04:05", v)
	return t
}
