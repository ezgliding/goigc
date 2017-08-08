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

// TaskType is the type of task.
type TaskType int

// Enum holding the possible TaskTypes.
const (
	TP1 TaskType = iota + 1
	TP2
	TP3
	TP4
	TP5
	FAI
)

// Result holds distance and the sequence of turn points.
type Result struct {
	Distance float64
	Points   []Point
}

// Optimizer implements an optimization algorithm for a track.
type Optimizer interface {
	Optimize(pts []Point, tp TaskType) (Result, error)
}
