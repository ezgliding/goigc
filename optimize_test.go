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
	"math"
	"testing"
)

const (
	errorMargin float64 = 0.02
)

type optimizeTest struct {
	name   string
	result map[int]float64
	margin float64
}

func (t *optimizeTest) validWithMargin(result float64, nTp int, errorMargin float64) bool {
	difference := math.Abs(t.result[nTp] - result)
	maxDifference := t.result[nTp] * errorMargin
	return difference < maxDifference
}

func (t *optimizeTest) valid(result float64, nTp int) bool {
	return t.validWithMargin(result, nTp, errorMargin)
}

var benchmarkTests = []optimizeTest{
	{
		name:   "optimize-short-flight-1",
		result: map[int]float64{1: 35.44619896425489},
	},
}

var optimizeTests = []optimizeTest{
	{
		name:   "optimize-short-flight-1",
		result: map[int]float64{1: 35.44619896425489, 2: 0.0, 3: 507.80108709626626},
	},
}

type distanceTest struct {
	t        string
	task     Task
	distance float64
}

var distanceTests = []distanceTest{
	{
		t: "all-points-the-same-distance-zero",
		task: Task{
			Start:      NewPointFromDMD("4453183N", "00512633E"),
			Turnpoints: []Point{NewPointFromDMD("4453183N", "00512633E")},
			Finish:     NewPointFromDMD("4453183N", "00512633E"),
		},
		distance: 0.0,
	},
	{
		t: "valid-task-sequence",
		task: Task{
			Start: NewPointFromDMD("4453183N", "00512633E"),
			Turnpoints: []Point{
				NewPointFromDMD("4353800N", "00615200E"),
				NewPointFromDMD("4506750N", "00633950E"),
				NewPointFromDMD("4424783N", "00644500E"),
			},
			Finish: NewPointFromDMD("4505550N", "00502883E"),
		},
		distance: 507.80108709626626,
	},
}

func TestDistance(t *testing.T) {
	var result float64
	for _, test := range distanceTests {
		t.Run(test.t, func(t *testing.T) {
			result = Distance(test.task)
			if result != test.distance {
				t.Errorf("expected %v got %v", test.distance, result)
			}
		})
	}
}
