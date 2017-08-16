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
	"math/rand"
	"sort"
)

// NewSimAnnealingOptimizer ...
func NewSimAnnealingOptimizer() Optimizer {
	return NewSimAnnealingOptimizerParams(10000, 1, 0.003)
}

// NewSimAnnealingOptimizerParams returns a BruteForceOptimizer with the given characteristics.
//
func NewSimAnnealingOptimizerParams(startTemperature float64, minTemperature float64, alpha float64) Optimizer {
	return &simAnnealingOptimizer{
		currentTemperature: startTemperature,
		minTemperature:     minTemperature,
		alpha:              alpha,
	}
}

type simAnnealingOptimizer struct {
	score              Score
	currentTemperature float64
	minTemperature     float64
	alpha              float64
	track              Track
	nPoints            int
	currentPoints      []int
	currentTask        Task
}

func (sa *simAnnealingOptimizer) initialize(track Track, nPoints int, score Score) {
	sa.track = track
	sa.nPoints = nPoints
	sa.score = score

	sa.currentPoints = make([]int, sa.nPoints)
	for i := 0; i < sa.nPoints; i++ {
		sa.currentPoints[i] = rand.Intn(len(sa.track.Points))
	}
	sort.Ints(sa.currentPoints)

	sa.currentTask = Task{}
	sa.currentTask.Start = sa.track.Points[sa.currentPoints[0]]
	for i := 1; i < sa.nPoints-1; i++ {
		sa.currentTask.Turnpoints[i] = sa.track.Points[sa.currentPoints[i]]
	}
	sa.currentTask.Finish = sa.track.Points[sa.currentPoints[sa.nPoints-1]]
}

func (sa *simAnnealingOptimizer) neighbour() ([]int, Task) {
	newPoints := sa.currentPoints
	newTask := sa.currentTask
	var prev, next int

	pos := rand.Intn(sa.nPoints)
	if prev = 0; pos-1 > 0 {
		prev = pos - 1
	}
	if next = sa.nPoints - 1; pos+1 < sa.nPoints-1 {
		next = pos + 1
	}
	value := newPoints[prev] + rand.Intn(newPoints[next]-newPoints[prev])
	newPoints[pos] = value

	switch pos {
	case 0:
		newTask.Start = sa.track.Points[newPoints[pos]]
	case sa.nPoints - 1:
		newTask.Finish = sa.track.Points[newPoints[pos]]
	default:
		newTask.Turnpoints[pos-1] = sa.track.Points[newPoints[pos]]
	}

	return newPoints, newTask
}

func (sa *simAnnealingOptimizer) acceptanceProb(task Task) float64 {
	diff := sa.score(sa.currentTask) - sa.score(task)
	return math.E * (diff / sa.currentTemperature)
}

func (sa *simAnnealingOptimizer) Optimize(track Track, nPoints int, score Score) (Task, error) {
	var acceptanceProb float64
	var points []int
	var task Task

	sa.initialize(track, nPoints, score)

	// loop while the temperature is above min
	for sa.currentTemperature > sa.minTemperature {
		points, task = sa.neighbour()
		acceptanceProb = sa.acceptanceProb(task)
		if acceptanceProb > rand.Float64() {
			sa.currentPoints = points
			sa.currentTask = task
		}
		sa.currentTemperature = sa.currentTemperature * sa.alpha
	}
	return task, nil
}
