// Copyright ©2015 Ricardo Rocha <rocha.porto@gmail.com>
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

import "math/rand"

type montecarlo struct {
	Cycles    int
	MCCycles  int
	MCGuesses int
}

// NewMontecarlo returns a new Montecarlo optimizer instance.
func NewMontecarloOptimizer() Optimizer {
	return &montecarlo{Cycles: 10, MCCycles: 100000, MCGuesses: 0}
}

func (mc *montecarlo) Optimize(track []Point, tp TaskType) (Result, error) {
	res := Result{Distance: 0.0, Points: make([]Point, tp)}

	nTP := int(tp)
	for c := 0; c < mc.Cycles; c++ {
		var candidate = make([]int, nTP)
		var cdistance float64

		// start with uniform distribution (equal distance)
		for i := 0; i < nTP; i++ {
			candidate[i] = ((len(track) - 1) / (nTP - 1)) * i
		}

		// run montecarlo cycles
		var index, bwp, twp, nwp int
		for i := 0; i < mc.MCCycles; i++ {
			index = rand.Int() % nTP
			bwp = 0
			if index > 0 {
				bwp = candidate[index-1]
			}
			twp = len(track) - 1
			if index < nTP-1 {
				twp = candidate[index+1]
			}
			nwp = rand.Intn(twp-bwp) + bwp
			candidate[index] = nwp
			//sort.Sort(sort.IntSlice(candidate))
			cdistance = mc.distance(track, candidate)
			if cdistance > res.Distance {
				for j := 0; j < nTP; j++ {
					res.Points[j] = track[candidate[j]]
				}
				res.Distance = cdistance
			}
		}
	}

	return res, nil
}

func (mc *montecarlo) distance(track []Point, tps []int) float64 {
	distance := 0.0
	var v float64
	for i := 0; i < len(tps)-1; i++ {
		v = track[tps[i]].Distance(track[tps[i+1]])
		distance += v
	}
	return distance
}
