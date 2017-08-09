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

// TaskType represents one of the possible task types (turnpoints and angles).
type TaskType int

// Enum holding the possible TaskTypes.
const (
	// Task with a single turnpoint (out and return).
	TP1 TaskType = iota + 1
	// Task with two turnpoints (generic triangle).
	TP2
	// Task with three turnpoints.
	TP3
	// Task with four turnpoints.
	TP4
	// Task with five turnpoints.
	TP5
	// Task with two turnpoints, forming a triangle where the shortest leg is at least 28% of the total.
	FAITriangle
)

// Optimizer returns an optimal Task with a number of turnpoints for the given Task.
//
// The optimal Task is selected taking into account the score function.
// Available score functions include MaxDistance and MaxPoints, but it is
// possible to pass the Optimizer a custom function.
type Optimizer interface {
	Optimize(track Track, tp TaskType) (Task, error)
}
