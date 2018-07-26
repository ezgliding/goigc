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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type ParseTest struct {
	t string
	c string
	e bool
}

var update = flag.Bool("update", false, "update golden test data")

func runTest(t *testing.T, ok bool, in, out string) {

	data, err := ioutil.ReadFile(in)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := Parse(string(data))
	if err != nil && !ok {
		return
	} else if err != nil {
		t.Error(err)
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("%v :: %+v", err, result)
	}

	// update golden if flag is passed
	if *update {
		if err = ioutil.WriteFile(out, resultJSON, 0644); err != nil {
			t.Fatal(err)
		}
	}

	expectedJSON, err := ioutil.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}

	if string(resultJSON) != string(expectedJSON) {
		t.Errorf("expected\n%+v\ngot\n%+v", string(expectedJSON), string(resultJSON))
	}
}

func TestParse(t *testing.T) {
	// testdata file name format is testname.[1|0].igc
	match, err := filepath.Glob("testdata/parse-*.igc")
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range match {
		t.Run(in, func(t *testing.T) {
			parts := strings.Split(in, ".")
			ok, _ := strconv.ParseBool(parts[len(parts)-2])
			out := fmt.Sprintf("%v.golden", in)

			runTest(t, ok, in, out)
		})
	}
}

func TestParseLocationMissing(t *testing.T) {
	_, err := ParseLocation("does-not-exist")
	if err == nil {
		t.Errorf("no error returned for missing file")
	}
}

func TestParseLocationEmpty(t *testing.T) {
	_, err := ParseLocation("")
	if err == nil {
		t.Errorf("no error returned empty string location")
	}
}

func TestStripUpToMissing(t *testing.T) {
	s := "nocolonhere"
	r := stripUpTo(s, ":")
	if r != s {
		t.Errorf("expected %v got %v", s, r)
	}
}

// Parse a given file and get a Track object.
func Example_parselocation() {
	track, _ := ParseLocation("sample-flight.igc")

	fmt.Printf("Track Pilot: %v", track.Pilot)
	fmt.Printf("Track Points %v", len(track.Pilot))
}

// Parse directly flight contents and get a Track object.
func Example_parsecontent() {
	// We could pass here a string with the full contents in IGC format
	track, _ := Parse(`
AFLA001Some Additional Data
HFDTE010203
HFFXA500
HFPLTPilotincharge:EZ PILOT
	`)

	fmt.Printf("Track Pilot: %v", track.Pilot)
	fmt.Printf("Track Points %v", len(track.Pilot))
}

func BenchmarkParse(b *testing.B) {
	c, err := ioutil.ReadFile("testdata/parse-0-benchmark-0.igc")
	if err != nil {
		b.Errorf("failed to load sample flight :: %v", err)
	}
	content := string(c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(content)
	}
}

func getTrack(task Task) Track {
	f := NewTrack()
	f.Task = task
	return f
}
