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
	"time"
)

// NewSimAnnealingOptimizer ...
func NewSimAnnealingOptimizer() Optimizer {
	return NewSimAnnealingOptimizerParams(1000, 1, 0.03, time.Now().UTC().UnixNano())
}

// NewSimAnnealingOptimizerParams returns a BruteForceOptimizer with the given characteristics.
func NewSimAnnealingOptimizerParams(startTemperature float64, minTemperature float64,
	alpha float64, seed int64) Optimizer {
	rand.Seed(seed)
	return &simAnnealingOptimizer{
		CurrentTemperature: startTemperature,
		MinTemperature:     minTemperature,
		Alpha:              alpha,
	}
}

type simAnnealingOptimizer struct {
	score              Score
	CurrentTemperature float64
	MinTemperature     float64
	Alpha              float64
	track              *Track
	nPoints            int
	candidate          Candidate
	best               Candidate
}

func (sa *simAnnealingOptimizer) initialize(track *Track, nPoints int, score Score) {
	sa.track = track
	sa.nPoints = nPoints
	sa.score = score
	sa.candidate = NewCandidateRandom(nPoints, track)
	sa.best = Candidate(sa.candidate)
}

func (sa *simAnnealingOptimizer) neighbour() Candidate {
	return sa.candidate.Neighbour()
}

func (sa *simAnnealingOptimizer) acceptanceProb(task Task) float64 {
	diff := sa.score(task) - sa.score(sa.candidate.Task)
	if diff > 0 {
		return 1.0
	}
	return math.E * (diff / sa.CurrentTemperature)
}

func (sa *simAnnealingOptimizer) Optimize(track Track, nPoints int, score Score) (Candidate, error) {
	var acceptanceProb float64
	var candidate Candidate

	sa.initialize(&track, nPoints, score)

	// loop while the temperature is above min
	for sa.CurrentTemperature > sa.MinTemperature {
		candidate = sa.neighbour()
		candidate.Metadata["Temperature"] = sa.CurrentTemperature
		candidate.Metadata["Best"] = sa.best
		acceptanceProb = sa.acceptanceProb(candidate.Task)
		if acceptanceProb > rand.Float64() {
			sa.candidate = candidate
		}
		candidate.Score = sa.score(candidate.Task)
		if candidate.Score > sa.best.Score {
			sa.best = Candidate(candidate)
		}
		sa.CurrentTemperature *= (1 - sa.Alpha)
	}
	return sa.best, nil
}
