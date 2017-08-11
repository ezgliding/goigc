// Copyright ©2017 The ezgliding Authors.
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
)

type point struct {
	lat string
	lng string
}

type PointTest struct {
	p1 point
	p2 point
	d  float64
}

var dmdTests = []PointTest{
	{
		p1: point{lat: "", lng: ""},
		p2: point{lat: "", lng: ""},
		d:  0,
	},
	{
		p1: point{lat: "5110179N", lng: "00102644W"},
		p2: point{lat: "5110179N", lng: "00102644W"},
		d:  0,
	},
	{
		p1: point{lat: "5110179N", lng: "00102644W"},
		p2: point{lat: "5230147N", lng: "00017612W"},
		d:  156.91393060997657,
	},
}

func TestDistanceDMD(t *testing.T) {
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
