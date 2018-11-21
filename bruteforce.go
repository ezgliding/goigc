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
)

type OType int

const (
	O0 OType = iota
	O1       = iota
	O2       = iota
	O3       = iota
	O4       = iota
	O5       = iota
)

// NewBruteForceOptimizer returns a BruteForceOptimizer with the given characteristics.
func NewBruteForceOptimizer(otype OType, cache bool) Optimizer {
	return &bruteForceOptimizer{otype: otype, cache: cache}
}

type bruteForceOptimizer struct {
	otype OType
	cache bool
}

func (b *bruteForceOptimizer) Optimize(track Track, nPoints int, score Score) (Task, error) {
	switch nPoints {
	case 1:
		if b.otype == O0 {
			return b.optimize1_o0(track, score)
		} else if b.otype == O1 {
			return b.optimize1_o1(track, score)
		} else if b.otype == O2 {
			return b.optimize1_o2(track, score)
		} else if b.otype == O3 {
			return b.optimize1_o3(track, score)
		} else if b.otype == O4 {
			return b.optimize1_o4(track, score)
		} else if b.otype == O5 {
			return b.optimize1_o5(track, score)
		} else {
			return Task{}, fmt.Errorf("invalid %v optimizer", b.otype)
		}
	case 2:
		return b.optimize2(track, score)
	case 3:
		return b.optimize3(track, score)
	default:
		return Task{}, fmt.Errorf("%v turn points not supported by this optimizer", nPoints)
	}
}

func (b *bruteForceOptimizer) optimize1_o0(track Track, score Score) (Task, error) {

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
				distance = task.DistanceSlow()
				if distance > optimalDistance {
					optimalDistance = distance
					optimalTask = Task(task)
				}
			}
		}
	}
	return optimalTask, nil
}

func (b *bruteForceOptimizer) optimize1_o1(track Track, score Score) (Task, error) {

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

func (b *bruteForceOptimizer) optimize1_o2(track Track, score Score) (Task, error) {

	var optimalDistance float64
	var distance float64
	var task Task
	var optimalTask Task

	task.Turnpoints = []Point{Point{}}
	optimalTask.Turnpoints = []Point{Point{}}
	for i := 0; i < len(track.Points)-2; i++ {
		for j := i + 1; j < len(track.Points)-1; j++ {
			for z := j + 1; z < len(track.Points); z++ {
				task.Start = track.Points[i]
				task.Turnpoints[0] = track.Points[j]
				task.Finish = track.Points[z]
				distance = task.Distance()
				if distance > optimalDistance {
					optimalDistance = distance
					optimalTask.Start = track.Points[i]
					optimalTask.Turnpoints[0] = track.Points[j]
					optimalTask.Finish = track.Points[z]
				}
			}
		}
	}
	return optimalTask, nil
}

func (b *bruteForceOptimizer) optimize1_o3(track Track, score Score) (Task, error) {

	var distance float64
	var task Task
	var optimalDistance float64
	var optimalTask Task
	var aa int

	iterations := (len(track.Points) - 2) * (len(track.Points) - 1)
	d := make(chan float64, iterations)
	t := make(chan Task, iterations)

	for i := 0; i < len(track.Points)-2; i++ {
		for j := i + 1; j < len(track.Points)-1; j++ {
			aa += 1
			go func(i, j int) {
				var distance float64
				var task Task
				var optimalDistance float64
				var optimalTask Task

				task.Turnpoints = []Point{Point{}}
				optimalTask.Turnpoints = []Point{Point{}}
				for z := j + 1; z < len(track.Points); z++ {
					task.Start = track.Points[i]
					task.Turnpoints[0] = track.Points[j]
					task.Finish = track.Points[z]
					distance = task.Distance()
					if distance > optimalDistance {
						optimalDistance = distance
						optimalTask.Start = track.Points[i]
						optimalTask.Turnpoints[0] = track.Points[j]
						optimalTask.Finish = track.Points[z]
					}
				}
				d <- optimalDistance
				t <- optimalTask
			}(i, j)
		}
	}

	for i := 0; i < 203841; i++ {
		distance = <-d
		task = <-t
		if distance > optimalDistance {
			optimalDistance = distance
			optimalTask = task
		}
	}

	return optimalTask, nil
}

func (b *bruteForceOptimizer) optimize1_o4(track Track, score Score) (Task, error) {

	var distance float64
	var task Task
	var optimalDistance float64
	var optimalTask Task

	d := make(chan float64, len(track.Points)-2)
	t := make(chan Task, len(track.Points)-2)

	for i := 0; i < len(track.Points)-2; i++ {
		go func(i int) {
			var optimalDistance float64
			var distance float64
			var task Task
			var optimalTask Task

			task.Turnpoints = []Point{Point{}}
			optimalTask.Turnpoints = []Point{Point{}}
			for j := i + 1; j < len(track.Points)-1; j++ {
				for z := j + 1; z < len(track.Points); z++ {
					task.Start = track.Points[i]
					task.Turnpoints[0] = track.Points[j]
					task.Finish = track.Points[z]
					distance = task.Distance()
					if distance > optimalDistance {
						optimalDistance = distance
						optimalTask.Start = track.Points[i]
						optimalTask.Turnpoints[0] = track.Points[j]
						optimalTask.Finish = track.Points[z]
					}
				}
			}
			d <- optimalDistance
			t <- optimalTask
		}(i)
	}

	for i := 0; i < len(track.Points)-2; i++ {
		distance = <-d
		task = <-t
		if distance > optimalDistance {
			optimalDistance = distance
			optimalTask = task
		}
	}

	return optimalTask, nil
}

func (b *bruteForceOptimizer) optimize1_o5(track Track, score Score) (Task, error) {

	var distance float64
	var task Task
	var optimalDistance float64
	var optimalTask Task

	n := len(track.Points)
	d := make(chan float64, n-2)
	t := make(chan Task, n-2)

	cache := make([]float64, n*(n-1))
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			cache[i*n+j] = track.Points[i].Distance(track.Points[j])
		}
	}

	for i := 0; i < n-2; i++ {
		go func(i int) {
			var optimalDistance float64
			var distance float64
			var task Task
			var optimalTask Task

			task.Turnpoints = []Point{Point{}}
			optimalTask.Turnpoints = []Point{Point{}}
			for j := i + 1; j < n-1; j++ {
				for z := j + 1; z < n; z++ {
					distance = cache[(i*n)+j] + cache[(j*n)+z]
					if distance > optimalDistance {
						optimalDistance = distance
						optimalTask.Start = track.Points[i]
						optimalTask.Turnpoints[0] = track.Points[j]
						optimalTask.Finish = track.Points[z]
					}
				}
			}
			d <- optimalDistance
			t <- optimalTask
		}(i)
	}

	for i := 0; i < n-2; i++ {
		distance = <-d
		task = <-t
		if distance > optimalDistance {
			optimalDistance = distance
			optimalTask = task
		}
	}

	return optimalTask, nil
}

func (b *bruteForceOptimizer) optimize2(track Track, score Score) (Task, error) {
	var distance float64
	var task Task
	var optimalDistance float64
	var optimalTask Task

	n := len(track.Points)
	d := make(chan float64, n-3)
	t := make(chan Task, n-3)

	cache := make([]float64, n*(n-1))
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			cache[i*n+j] = track.Points[i].Distance(track.Points[j])
		}
	}

	for i := 0; i < n-3; i++ {
		go func(i int) {
			var optimalDistance float64
			var distance float64
			var task Task
			var optimalTask Task

			task.Turnpoints = make([]Point, 2)
			optimalTask.Turnpoints = make([]Point, 2)
			for j := i + 1; j < n-2; j++ {
				for w := j + 1; w < n-1; w++ {
					for z := w + 1; z < n; z++ {
						distance = cache[(i*n)+j] + cache[(j*n)+w] + cache[(w*n)+z]
						if distance > optimalDistance {
							optimalDistance = distance
							optimalTask.Start = track.Points[i]
							optimalTask.Turnpoints[0] = track.Points[j]
							optimalTask.Turnpoints[1] = track.Points[w]
							optimalTask.Finish = track.Points[z]
						}
					}
				}
			}
			d <- optimalDistance
			t <- optimalTask
		}(i)
	}

	for i := 0; i < n-3; i++ {
		distance = <-d
		task = <-t
		if distance > optimalDistance {
			optimalDistance = distance
			optimalTask = task
		}
	}

	return optimalTask, nil
}

func (b *bruteForceOptimizer) optimize3(track Track, score Score) (Task, error) {
	var distance float64
	var task Task
	var optimalDistance float64
	var optimalTask Task

	n := len(track.Points)
	d := make(chan float64, n-3)
	t := make(chan Task, n-3)

	cache := make([]float64, n*(n-1))
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			cache[i*n+j] = track.Points[i].Distance(track.Points[j])
		}
	}

	for i := 0; i < n-3; i++ {
		go func(i int) {
			var optimalDistance float64
			var distance float64
			var task Task
			var optimalTask Task

			task.Turnpoints = make([]Point, 3)
			optimalTask.Turnpoints = make([]Point, 3)
			for j := i + 1; j < n-2; j++ {
				for k := j + 1; k < n-1; k++ {
					for w := k + 1; w < n; w++ {
						for z := w + 1; z < n; z++ {
							distance = cache[(i*n)+j] + cache[(j*n)+k] + cache[(k*n)+w] + cache[(w*n)+z]
							if distance > optimalDistance {
								optimalDistance = distance
								optimalTask.Start = track.Points[i]
								optimalTask.Turnpoints[0] = track.Points[j]
								optimalTask.Turnpoints[1] = track.Points[k]
								optimalTask.Turnpoints[2] = track.Points[w]
								optimalTask.Finish = track.Points[z]
							}
						}
					}
				}
			}
			d <- optimalDistance
			t <- optimalTask
		}(i)
	}

	for i := 0; i < n-3; i++ {
		distance = <-d
		task = <-t
		if distance > optimalDistance {
			optimalDistance = distance
			optimalTask = task
		}
	}

	return optimalTask, nil
}
