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
