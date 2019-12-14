// Copyright The ezgliding Authors.
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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type phaseTest struct {
	name   string
	result string //nolint
}

var phaseTests = []phaseTest{
	{
		name:   "phases-short-flight-1",
		result: "phases-short-flight-1.simple",
	},
	{
		name:   "phases-long-flight-1",
		result: "phases-long-flight-1.simple",
	},
}

func TestPhases(t *testing.T) {
	for _, test := range phaseTests {
		t.Run(fmt.Sprintf("%v\n", test.name), func(t *testing.T) {
			f := filepath.Join("../../testdata/phases", fmt.Sprintf("%v.igc", test.name))
			golden := fmt.Sprintf("../../testdata/phases/%v.golden.igc", test.name)
			track, err := ParseLocation(f)
			if err != nil {
				t.Fatal(err)
			}
			phases, err := track.Phases()
			if err != nil {
				t.Fatal(err)
			}
			for _, ph := range phases {
				fmt.Printf("%v %v %v\n", ph.Type, ph.Start.Time, ph.End.Time)
			}

			// update golden if flag is passed
			if *update {
				jsn, err := json.Marshal(phases)
				if err != nil {
					t.Fatalf("%+v :: %v\n", phases, err)
				}
				if err = ioutil.WriteFile(golden, jsn, 0644); err != nil {
					t.Fatal(err)
				}
			}

			b, _ := ioutil.ReadFile(golden)
			var goldenPhases []Phase
			_ = json.Unmarshal(b, &goldenPhases)
			if len(phases) != len(goldenPhases) {
				t.Errorf("expected %v got %v phases", len(goldenPhases), len(phases))
			}
		})
	}
}

func TestEncodePhasesKML(t *testing.T) {
	for _, test := range phaseTests {
		t.Run(fmt.Sprintf("%v\n", test.name), func(t *testing.T) {
			f := filepath.Join("../../testdata/phases", fmt.Sprintf("%v.igc", test.name))
			golden := fmt.Sprintf("%v.golden.kml", f)
			track, err := ParseLocation(f)
			if err != nil {
				t.Fatal(err)
			}

			kml, err := track.encodePhasesKML()
			if err != nil {
				t.Fatal(err)
			}

			buf := new(bytes.Buffer)
			err = kml.WriteIndent(buf, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			// update golden if flag is passed
			if *update {
				if err = ioutil.WriteFile(golden, buf.Bytes(), 0644); err != nil {
					t.Fatal(err)
				}
			}

			b, _ := ioutil.ReadFile(golden)
			if string(b) != buf.String() {
				t.Errorf("expected\n%v\ngot\n%v\n", string(b), buf.String())
			}
		})
	}
}
