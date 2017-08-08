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

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type ParseTest struct {
	t string
	c string
	e bool
}

var parseTests = []ParseTest{
	{
		"parse-basic-header-test",
		`
AFLA001Some Additional Data
HFDTE010203
HFFXA500
HFPLTPilotincharge:EZ PILOT
HFCM2Crew2:EZ CREW
HFGTYGliderType:EZ TYPE
HFGIDGliderID:EZ ID
HFDTM100GPSDatum:WGS84
HFRFWFirmwareVersion:v 0.1
HFRHWHardwareVersion:v 0.2
HFFTYFRType:EZ RECORDER,001
HFGPSEZ GPS,002,12,5000
HFPRSPressAltSensor:EZ PRESSURE
HFCIDCompetitionID:EZ COMPID
HFCCLCompetitionClass:EZ COMPCLASS
HFTZNTimezone:2.00
`,
		false,
	},
	{"parse-A-record-failure-too-short",
		"AFLA0", true},
	{"parse-H-record-failure-too-short",
		"HFFX", true},
	{"parse-H-record-failure-bad-date",
		"HFDTE330203", true},
	{"parse-H-record-failure-date-too-short",
		"HFDTE33", true},
	{"parse-H-record-failure-bad-fix-accuracy",
		"HFFXAAAA", true},
	{"parse-H-record-failure-fix-accuracy-too-short",
		"HFFXA20", true},
	{"parse-H-record-failure-gps-datum-too-short",
		"HFDTM20", true},
	{"parse-H-record-failure-unknown-field",
		"HFZZZaaa", true},
	{"parse-H-record-failure-bad-timezone",
		"HFTZNaa", true},
	{
		"parse-basic-flight-test",
		`
I033638FXA3940SIU4143ENL
J010812HDT
C150701213841160701000102500KTri
C5111359N00101899WEZ TAKEOFF
C5110179N00102644WEZ START
C5209092N00255227WEZ TP1
C5230147N00017612WEZ TP2
C5110179N00102644WEZ FINISH
C5111359N00101899WEZ LANDING
F160240040609123624
D20331
E160245ATS102312
B1602455107126N00149300WA002880042919509020
K16024800090
B1603105107212N00149174WV002930043519608024
LPLTLOG TEXT
GREJNGJERJKNJKRE31895478537H43982FJN9248F942389T433T
GJNJK2489IERGNV3089IVJE9GO398535J3894N358954983O0934
`,
		false,
	},
	{"parse-point-fix-wrong-size",
		"B110001", true},
	{"parse-point-fix-bad-time",
		"B3103105107212N00149174WV002930043519608024", true},
	{"parse-point-fix-bad-fix-validity",
		"B1603105107212N00149174WX002930043519608024", true},
	{"parse-point-fix-bad-pressure-altitude",
		"B1603105107212N00149174WV0029a0043519608024", true},
	{"parse-point-fix-bad-gnss-altitude",
		"B1603105107212N00149174WV002930043a19608024", true},
	{"parse-irecord-wrong-size",
		"I0", true},
	{"parse-irecord-invalid-value-for-field-number",
		"I0a", true},
	{"parse-irecord-wrong-size-with-fields",
		"I02AAA0102BBB030", true},
	{"parse-jrecord-wrong-size",
		"J0", true},
	{"parse-jrecord-invalid-value-for-field-number",
		"J0a", true},
	{"parse-jrecord-wrong-size-with-fields",
		"J02AAA0102BBB030", true},
	{"parse-k-wrong-size",
		"K16024", true},
	{"parse-k-invalid-date",
		"K160271", true},
	{"parse-k-wrong-size",
		"K16027000090", true},
	{"parse-e-wrong-size",
		"E16024", true},
	{"parse-e-invalid-date",
		"E160271ATS", true},
	{"parse-f-wrong-size",
		"F16024", true},
	{"parse-f-invalid-date",
		"F1602710102", true},
	{"parse-f-invalid-num-satellites",
		"F1602310a02", true},
	{"parse-l-wrong-size",
		"LPL", true},
	{"parse-c-bad-num-lines",
		"C150701213841160701000102500KTri", true},
	{"parse-c-wrong-size-first-line",
		"C15070121384116070100010", true},
	{"parse-c-invalid-num-of-tps",
		"C15070121384116070100010a", true},
	{"parse-c-invalid-declaration-date",
		"C350701213841160701000102500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5230147N00017612WEZ TP2\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", false},
	{"parse-c-invalid-flight-date",
		"C150701213841360701000102500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5230147N00017612WEZ TP2\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", false},
	{"parse-c-invalid-task-number",
		"C150701213841160701000a01500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", true},
	{"parse-c-invalid-takeoff",
		"C150701213841160701000101500KTri\nC5111359N00101899\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", true},
	{"parse-c-invalid-start",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", true},
	{"parse-c-invalid-tp",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", true},
	{"parse-c-invalid-finish",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644\nC5111359N00101899WEZ LANDING", true},
	{"parse-c-invalid-landing",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899", true},
	{"parse-d-wrong-size",
		"D2033", true},
	{"parse-invalid-record",
		"RANDOM GARBAGE", true},
}

var update = flag.Bool("update", false, "update golden test data")

func Get(t *testing.T, actual Track, test string) Track {
	golden := filepath.Join("tdata", fmt.Sprintf("%s.json", test))

	actualJson, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatalf("%v :: %+v", err, actual)
	}

	if *update {
		if err = ioutil.WriteFile(golden, actualJson, 0644); err != nil {
			t.Fatal(err)
		}
	}

	expectedJson, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	var expected Track
	err = json.Unmarshal(expectedJson, &expected)
	if err != nil {
		t.Fatal(err)
	}

	return expected
}

func TestParse(t *testing.T) {
	var expected Track
	for _, test := range parseTests {
		result, err := Parse(test.c)
		if err != nil && test.e {
			continue
		} else if err != nil {
			t.Errorf("%v failed :: %v", test.t, err)
			continue
		}
		expected = Get(t, result, test.t)
		resultJson, _ := json.Marshal(result)
		expectedJson, _ := json.Marshal(expected)
		if string(resultJson) != string(expectedJson) {
			t.Errorf("%v failed :: expected\n%+v\ngot\n%+v", test.t, expected, result)
			continue
		}
	}
}

func TestStripUpToMissing(t *testing.T) {
	s := "nocolonhere"
	r := stripUpTo(s, ":")
	if r != s {
		t.Errorf("expected %v got %v", s, r)
	}
}

func BenchmarkParse(b *testing.B) {
	c, err := ioutil.ReadFile("../test/sample-flight.igc")
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
