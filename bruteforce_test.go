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
	"path/filepath"
	"testing"
)

func TestBruteForceOptimize(t *testing.T) {
	opt := NewBruteForceOptimizer(O0, false)

	for _, test := range optimizeTests {
		for tp, expected := range test.result {
			if tp > 1 {
				continue
			}
			t.Run(fmt.Sprintf("%v/%v", test.name, tp), func(t *testing.T) {
				track, err := ParseLocation(filepath.Join("testdata", fmt.Sprintf("%v.igc", test.name)))
				if err != nil {
					t.Fatal(err)
				}
				task, err := opt.Optimize(track, tp, Distance)
				if err != nil {
					t.Fatal(err)
				}
				result := task.Distance()
				if !test.valid(result, tp) {
					t.Errorf("expected %v got %v", expected, result)
				}
			})
		}
	}
}

func BenchmarkBruteForceOptimizeO0(b *testing.B) {
	benchmarkBruteForce(O0, b)
}

func BenchmarkBruteForceOptimizeO1(b *testing.B) {
	benchmarkBruteForce(O1, b)
}

func BenchmarkBruteForceOptimizeO2(b *testing.B) {
	benchmarkBruteForce(O2, b)
}

func BenchmarkBruteForceOptimizeO3(b *testing.B) {
	benchmarkBruteForce(O3, b)
}

func BenchmarkBruteForceOptimizeO4(b *testing.B) {
	benchmarkBruteForce(O4, b)
}

func BenchmarkBruteForceOptimizeO5(b *testing.B) {
	benchmarkBruteForce(O5, b)
}

func benchmarkBruteForce(otype OType, b *testing.B) {
	opt := NewBruteForceOptimizer(otype, false)

	for _, test := range benchmarkTests {
		for tp, expected := range test.result {
			if tp > 1 && otype != O5 {
				continue
			}
			track, err := ParseLocation(filepath.Join("testdata", fmt.Sprintf("%v.igc", test.name)))
			if err != nil {
				b.Fatal(err)
			}
			b.Run(fmt.Sprintf("%v/%v", test.name, tp), func(b *testing.B) {
				task, err := opt.Optimize(track, tp, Distance)
				if err != nil {
					b.Fatal(err)
				}
				result := task.Distance()
				if !test.valid(result, tp) {
					b.Errorf("expected %v got %v", expected, result)
				}
			})
		}
	}
}
