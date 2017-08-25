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
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	return math.E * (diff / sa.currentTemperature)
}

type RunLog struct {
	Candidates  []Candidate
	Bests       []Candidate
	Temperature []float64
	Points      []Point
}

func (sa *simAnnealingOptimizer) Optimize(track Track, nPoints int, score Score) (Task, error) {
	var acceptanceProb float64
	var candidate Candidate

	sa.initialize(&track, nPoints, score)

	log := RunLog{Candidates: make([]Candidate, 1), Bests: make([]Candidate, 1), Temperature: make([]float64, 1), Points: track.Points}
	// loop while the temperature is above min
	for sa.currentTemperature > sa.minTemperature {
		candidate = sa.neighbour()
		acceptanceProb = sa.acceptanceProb(candidate.Task)
		if acceptanceProb > rand.Float64() {
			sa.candidate = candidate
		}
		candidate.Score = sa.score(candidate.Task)
		if candidate.Score > sa.best.Score {
			sa.best = Candidate(candidate)
		}
		sa.currentTemperature *= (1 - sa.alpha)
		log.Candidates = append(log.Candidates, candidate)
		log.Bests = append(log.Bests, sa.best)
		log.Temperature = append(log.Temperature, sa.currentTemperature)
	}
	txt, _ := json.Marshal(log)
	ioutil.WriteFile("tmpdata.js", []byte(fmt.Sprintf("var data = '%v';", string(txt))), 0644)
	return sa.best.Task, nil
}
