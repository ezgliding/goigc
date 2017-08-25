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
	"math/rand"
	"sort"
)

// Score functions calculate a score for the given Task.
//
// The main use of these functions is in passing them to the Optimizers, so
// they can evaluate each Task towards different goals.
//
// Example functions include the total distance between all turn points or an
// online competition (netcoupe, online contest) score which takes additional
// metrics of each leg into account.
type Score func(task Task) float64

// Distance returns the sum of distances between each of the points in the Task.
//
// The sum is made calculating the distances between each two consecutive Points.
func Distance(task Task) float64 {
	return task.Distance()
}

// Optimizer returns an optimal Task for the given turnpoints and Score function.
//
// Available score functions include MaxDistance and MaxPoints, but it is
// possible to pass the Optimizer a custom function.
//
// Optimizers might not support a high number of turnpoints. As an example, the
// BruteForceOptimizer does not perform well with nPoints > 2, and might decide
// to return an error instead of attempting to finalize the optimization
// indefinitely.
type Optimizer interface {
	Optimize(track Track, nPoints int, score Score) (Task, error)
}

// Candidate aggregates Point indexes and a corresponding Task.
//
// It is a utility for Optimizer implementations requiring the tracking of the
// Point indexes to calculate neighbours and variations of a current candidate.
type Candidate struct {
	indexes []int
	track   *Track
	Task    Task
	Score   float64
}

// NewCandidate returns a Candidate initialized with nPoints, all set to zero.
func NewCandidate(nPoints int, track *Track) Candidate {
	return Candidate{indexes: make([]int, nPoints), track: track, Task: NewTask(nPoints)}
}

// NewCandidateRandom returns a Candidate with a random set of nPoints of track.
func NewCandidateRandom(nPoints int, track *Track) Candidate {
	candidate := NewCandidate(nPoints, track)
	for i := 0; i < nPoints; i++ {
		candidate.indexes[i] = rand.Intn(len(track.Points))
	}
	sort.Ints(candidate.indexes)
	for i := 0; i < nPoints; i++ {
		candidate.Task.Turnpoints[i] = track.Points[candidate.indexes[i]]
	}
	return candidate
}

// Neighbour returns a neighbouring Candidate to this one.
//
// The selection is done by randomly choosing one of the turn points (p), and
// changing that point (p) with another on the Track, where the new point is
// between p-1 and p+1.
func (c *Candidate) Neighbour() Candidate {
	var value int
	nPoints := len(c.indexes)

	candidate := NewCandidate(nPoints, c.track)
	pos := rand.Intn(nPoints)
	switch pos {
	case 0:
		value = randMinMax(0, c.indexes[pos]-1)
	case nPoints - 1:
		value = randMinMax(c.indexes[pos]+1, len(c.track.Points)-1)
	default:
		value = randMinMax(c.indexes[pos-1]+1, c.indexes[pos+1]-1)
	}

	for i := 0; i < nPoints; i++ {
		candidate.indexes[i] = c.indexes[i]
		candidate.Task.Turnpoints[i] = c.track.Points[candidate.indexes[i]]
	}
	candidate.indexes[pos] = value
	candidate.Task.Turnpoints[pos] = c.track.Points[value]

	return candidate
}

func randMinMax(min int, max int) int {
	if max-min < 0 {
		return 0
	} else if max-min == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}
