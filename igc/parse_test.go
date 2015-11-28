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
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/kellydunn/golang-geo"
)

type ParseTest struct {
	t string
	c string
	r Track
	e bool
}

var parseTests = []ParseTest{
	{
		"basic header test",
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
		Track{
			Header: Header{
				Manufacturer: "FLA", UniqueID: "001", AdditionalData: "Some Additional Data",
				Date:        time.Date(2003, time.February, 01, 0, 0, 0, 0, time.UTC),
				FixAccuracy: 500, Pilot: "EZ PILOT", Crew: "EZ CREW",
				GliderType: "EZ TYPE", GliderID: "EZ ID", GPSDatum: "WGS84",
				FirmwareVersion: "v 0.1", HardwareVersion: "v 0.2",
				FlightRecorder: "EZ RECORDER,001", GPS: "EZ GPS,002,12,5000",
				PressureSensor: "EZ PRESSURE", CompetitionID: "EZ COMPID",
				CompetitionClass: "EZ COMPCLASS", Timezone: *time.FixedZone("", 2*3600),
			},
			K:          map[time.Time]map[string]string{},
			Events:     map[time.Time]map[string]string{},
			Satellites: map[time.Time][]int{},
		},
		false,
	},
	{"A record failure too short",
		"AFLA0", Track{}, true},
	{"H record failure too short",
		"HFFX", Track{}, true},
	{"H record failure bad date",
		"HFDTE330203", Track{}, true},
	{"H record failure date too short",
		"HFDTE33", Track{}, true},
	{"H record failure bad fix accuracy",
		"HFFXAAAA", Track{}, true},
	{"H record failure fix accuracy too short",
		"HFFXA20", Track{}, true},
	{"H record failure gps datum too short",
		"HFDTM20", Track{}, true},
	{"H record failure unknown field",
		"HFZZZaaa", Track{}, true},
	{"H record failure bad timezone",
		"HFTZNaa", Track{}, true},
	{
		"basic flight test",
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
		Track{
			Points: []Point{
				Point{
					Point:            *geo.NewPoint(51.118766666666666, -1.8216666666666668),
					Time:             time.Date(0, 1, 1, 16, 2, 45, 0, time.UTC),
					FixValidity:      'A',
					PressureAltitude: 288, GNSSAltitude: 429,
					IData: map[string]string{
						"FXA": "195", "SIU": "09", "ENL": "020",
					},
					NumSatellites: 6,
				},
				Point{
					Point:       *geo.NewPoint(51.118766666666666, -1.8216666666666668),
					Time:        time.Date(0, 1, 1, 16, 3, 10, 0, time.UTC),
					FixValidity: 'V', PressureAltitude: 293, GNSSAltitude: 435,
					IData: map[string]string{
						"FXA": "196", "SIU": "08", "ENL": "024",
					},
					NumSatellites: 6,
				},
			},
			K: map[time.Time]map[string]string{
				time.Date(0, 1, 1, 16, 2, 48, 0, time.UTC): map[string]string{
					"HDT": "00090",
				},
			},
			Events: map[time.Time]map[string]string{
				time.Date(0, 1, 1, 16, 2, 45, 0, time.UTC): map[string]string{
					"ATS": "102312",
				},
			},
			Satellites: map[time.Time][]int{
				time.Date(0, 1, 1, 16, 02, 40, 0, time.UTC): []int{4, 6, 9, 12, 36, 24},
			},
			Logbook: []LogEntry{
				LogEntry{Type: "PLT", Text: "LOG TEXT"},
			},
			Task: Task{
				DeclarationDate: time.Date(2001, time.July, 15, 21, 38, 41, 0, time.UTC),
				Date:            time.Date(2001, time.July, 16, 0, 0, 0, 0, time.UTC),
				Number:          1,
				Takeoff: Point{
					Point:       *geo.NewPoint(51.18931666666667, -1.03165),
					Description: "EZ TAKEOFF"},
				Start: Point{
					Point:       *geo.NewPoint(51.16965, -1.0440666666666667),
					Description: "EZ START"},
				Turnpoints: []Point{
					Point{
						Point:       *geo.NewPoint(52.15153333333333, -2.9204499999999998),
						Description: "EZ TP1"},
					Point{
						Point:       *geo.NewPoint(52.50245, -0.2935333333333333),
						Description: "EZ TP2"},
				},
				Finish: Point{
					Point:       *geo.NewPoint(51.16965, -1.0440666666666667),
					Description: "EZ FINISH"},
				Landing: Point{
					Point:       *geo.NewPoint(51.18931666666667, -1.03165),
					Description: "EZ LANDING"},
				Description: "500KTri",
			},
			DGPSStationID: "0331",
			Signature:     "REJNGJERJKNJKRE31895478537H43982FJN9248F942389T433TJNJK2489IERGNV3089IVJE9GO398535J3894N358954983O0934",
		},
		false,
	},
	{"point/fix wrong size",
		"B110001", Track{}, true},
	{"point/fix bad time",
		"B3103105107212N00149174WV002930043519608024", Track{}, true},
	{"point/fix bad fix validity",
		"B1603105107212N00149174WX002930043519608024", Track{}, true},
	{"point/fix bad pressure altitude",
		"B1603105107212N00149174WV0029a0043519608024", Track{}, true},
	{"point/fix bad gnss altitude",
		"B1603105107212N00149174WV002930043a19608024", Track{}, true},
	{"irecord wrong size",
		"I0", Track{}, true},
	{"irecord invalid value for field number",
		"I0a", Track{}, true},
	{"irecord wrong size with fields",
		"I02AAA0102BBB030", Track{}, true},
	{"jrecord wrong size",
		"J0", Track{}, true},
	{"jrecord invalid value for field number",
		"J0a", Track{}, true},
	{"jrecord wrong size with fields",
		"J02AAA0102BBB030", Track{}, true},
	{"k wrong size",
		"K16024", Track{}, true},
	{"k invalid date",
		"K160271", Track{}, true},
	{"k wrong size",
		"K16027000090", Track{}, true},
	{"e wrong size",
		"E16024", Track{}, true},
	{"e invalid date",
		"E160271ATS", Track{}, true},
	{"f wrong size",
		"F16024", Track{}, true},
	{"f invalid date",
		"F1602710102", Track{}, true},
	{"f invalid num satellites",
		"F1602310a02", Track{}, true},
	{"l wrong size",
		"LPL", Track{}, true},
	{"c bad num lines",
		"C150701213841160701000102500KTri", Track{}, true},
	{"c wrong size first line",
		"C15070121384116070100010", Track{}, true},
	{"c invalid num of tps",
		"C15070121384116070100010a", Track{}, true},
	{"c invalid declaration date",
		"C350701213841160701000102500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5230147N00017612WEZ TP2\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", getTrack(Task{
			DeclarationDate: time.Time{},
			Date:            time.Date(2001, time.July, 16, 0, 0, 0, 0, time.UTC),
			Number:          1,
			Takeoff: Point{
				Point:       *geo.NewPoint(51.18931666666667, -1.03165),
				Description: "EZ TAKEOFF"},
			Start: Point{
				Point:       *geo.NewPoint(51.16965, -1.0440666666666667),
				Description: "EZ START"},
			Turnpoints: []Point{
				Point{
					Point:       *geo.NewPoint(52.15153333333333, -2.9204499999999998),
					Description: "EZ TP1"},
				Point{
					Point:       *geo.NewPoint(52.50245, -0.2935333333333333),
					Description: "EZ TP2"},
			},
			Finish: Point{
				Point:       *geo.NewPoint(51.16965, -1.0440666666666667),
				Description: "EZ FINISH"},
			Landing: Point{
				Point:       *geo.NewPoint(51.18931666666667, -1.03165),
				Description: "EZ LANDING"},
			Description: "500KTri",
		}), false},
	{"c invalid flight date",
		"C150701213841360701000102500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5230147N00017612WEZ TP2\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", getTrack(Task{
			DeclarationDate: time.Date(2001, time.July, 15, 21, 38, 41, 0, time.UTC),
			Date:            time.Time{},
			Number:          1,
			Takeoff: Point{
				Point:       *geo.NewPoint(51.18931666666667, -1.03165),
				Description: "EZ TAKEOFF"},
			Start: Point{
				Point:       *geo.NewPoint(51.16965, -1.0440666666666667),
				Description: "EZ START"},
			Turnpoints: []Point{
				Point{
					Point:       *geo.NewPoint(52.15153333333333, -2.9204499999999998),
					Description: "EZ TP1"},
				Point{
					Point:       *geo.NewPoint(52.50245, -0.2935333333333333),
					Description: "EZ TP2"},
			},
			Finish: Point{
				Point:       *geo.NewPoint(51.16965, -1.0440666666666667),
				Description: "EZ FINISH"},
			Landing: Point{
				Point:       *geo.NewPoint(51.18931666666667, -1.03165),
				Description: "EZ LANDING"},
			Description: "500KTri",
		}), false},
	{"c invalid task number",
		"C150701213841160701000a01500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", Track{}, true},
	{"c invalid takeoff",
		"C150701213841160701000101500KTri\nC5111359N00101899\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", Track{}, true},
	{"c invalid start",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", Track{}, true},
	{"c invalid tp",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227\nC5110179N00102644WEZ FINISH\nC5111359N00101899WEZ LANDING", Track{}, true},
	{"c invalid finish",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644\nC5111359N00101899WEZ LANDING", Track{}, true},
	{"c invalid landing",
		"C150701213841160701000101500KTri\nC5111359N00101899WEZ TAKEOFF\nC5110179N00102644WEZ START\nC5209092N00255227WEZ TP1\nC5110179N00102644WEZ FINISH\nC5111359N00101899", Track{}, true},
	{"d wrong size",
		"D2033", Track{}, true},
	{"invalid record",
		"RANDOM GARBAGE", Track{}, true},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		result, err := Parse(test.c)
		if err != nil && test.e {
			continue
		} else if err != nil {
			t.Errorf("%v failed :: %v", test.t, err)
			continue
		}
		if !reflect.DeepEqual(result, test.r) {
			t.Errorf("%v failed :: expected\n%+v\ngot\n%+v", test.t, test.r, result)
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
	c, err := ioutil.ReadFile("t/sample-igc")
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
