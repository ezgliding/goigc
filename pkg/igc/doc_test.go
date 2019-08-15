// Copyright 2017 The ezgliding Authors.
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
)

// Parse and ParseLocation return a Track object.
func Example_parse() {
	// We can parse passing a file location
	track, _ := ParseLocation("sample-flight.igc")

	// Or similarly giving directly the contents
	contents := `
AFLA001Some Additional Data
HFDTE010203
HFFXA500
HFPLTPilotincharge:EZ PILOT
	`
	track, _ = Parse(contents)

	// Accessing track metadata
	fmt.Printf("Track Pilot: %v", track.Pilot)
	fmt.Printf("Track Points %v", len(track.Pilot))
}

// Calculate the total track distance using the Points.
//
// This is not a very useful metric (you should look at one of the Optimizers)
// instead, but it is a good example of how to use the Point data in the Track.
func Example_totaldistance() {
	track, _ := ParseLocation("sample-flight.igc")
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}
	fmt.Printf("Distance was %v", totalDistance)
}

// Calculate the optimal track distance for the multiple possible tasks using the Brute Force Optimizer:
func Example_optimize() {
	track, _ := ParseLocation("sample-flight.igc")

	// In this case we use a brute force optimizer
	o := NewBruteForceOptimizer(false)
	r, _ := o.Optimize(track, 1, Distance)
	fmt.Printf("Optimization result was: %v", r)
}
