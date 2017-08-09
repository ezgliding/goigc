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
	"fmt"
	"time"
)

// NewBruteForceOptimizer returns a BruteForceOptimizer with the given characteristics.
//
func NewBruteForceOptimizer(cache bool) Optimizer {
	return &bruteForceOptimizer{cache: cache}
}

type bruteForceOptimizer struct {
	cache bool
}

func (b *bruteForceOptimizer) Optimize(track Track, nPoints int, score Score) (Task, error) {
	time.Sleep(5 * time.Second)
	switch nPoints {

	case 1:
		return b.optimize1(track, score)
	case 2:
		return b.optimize2(track, score)
	default:
		return Task{}, fmt.Errorf("%v turn points not supported by this optimizer", nPoints)
	}
}

func (b *bruteForceOptimizer) optimize1(track Track, score Score) (Task, error) {

	var optimalDistance float64
	var distance float64
	var task Task
	var optimalTask Task

	for i := 0; i < len(track.Points)-2; i++ {
		for j := i + 1; j < len(track.Points)-1; j++ {
			for z := j + 1; z < len(track.Points); z++ {
				task = Task{
					Start:      track.Points[i],
					Turnpoints: []Point{track.Points[j]},
					Finish:     track.Points[z],
				}
				distance = task.Distance()
				if distance > optimalDistance {
					optimalDistance = distance
					optimalTask = Task(task)
				}
			}
		}
	}
	return optimalTask, nil
}

func (b *bruteForceOptimizer) optimize2(track Track, score Score) (Task, error) {

	var optimalDistance float64
	var distance float64
	var optimalTask Task

	for i := 0; i < len(track.Points)-3; i++ {
		for j := i + 1; j < len(track.Points)-2; j++ {
			for w := j + 1; w < len(track.Points)-1; w++ {
				for z := w + 1; z < len(track.Points); z++ {
					task := Task{
						Start:      track.Points[i],
						Turnpoints: []Point{track.Points[j], track.Points[w]},
						Finish:     track.Points[z],
					}
					distance = task.Distance()
					if distance > optimalDistance {
						optimalDistance = distance
						optimalTask = task
					}
				}
			}
		}
	}
	return optimalTask, nil
}
