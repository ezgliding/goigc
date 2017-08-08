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

// Scorer is the interface that wraps the Score method.
//
// Scores takes a given set of points and returns its score possibly taking
// into account the given task.
type Scorer interface {
	Score(pts []Point, t Task) (float64, error)
}

// MaxDistance ...
type MaxDistance string

// Score ...
func (md MaxDistance) Score(pts []Point, t Task) (float64, error) {
	d := 0.0
	for i := 0; i < len(pts)-1; i++ {
		d += pts[i].Distance(pts[i+1])
	}
	return d, nil
}
