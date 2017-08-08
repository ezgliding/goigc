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

import "fmt"

// NewBruteForceOptimizer ...
func NewBruteForceOptimizer(cache bool) Optimizer {
	return &bruteForceOptimizer{cache: cache}
}

type bruteForceOptimizer struct {
	cache bool
}

func (b *bruteForceOptimizer) Optimize(pts []Point, t TaskType) (Result, error) {
	res := Result{Distance: 0.0, Points: make([]Point, 3)}

	var i, j, z, cnt int
	var d float64
	switch t {
	case TP1:
		for i = 0; i < len(pts)-2; i++ {
			for j = i + 1; j < len(pts)-1; j++ {
				for z = j + 1; z < len(pts); z++ {
					d = pts[i].Distance(pts[j])
					d += pts[j].Distance(pts[z])
					if d > res.Distance {
						res.Distance = d
						res.Points[0] = pts[i]
						res.Points[1] = pts[j]
						res.Points[2] = pts[z]
					}
					cnt++
				}
			}
		}
	default:
		return res, fmt.Errorf("unsupported task type (%v)", t)
	}
	return res, nil
}
