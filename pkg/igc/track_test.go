// Copyright ©2019 The ezgliding Authors.
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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/geo/s2"
)

const (
	testDir = "/tmp/ezgliding/tests"
)

type simplifyTest struct {
	name   string
	result string
}

var simplifyTests = []simplifyTest{
	{
		name:   "simplify-short-flight-1",
		result: "simplify-short-flight-1.simple",
	},
}

func TestSimplify(t *testing.T) {
	for _, test := range simplifyTests {
		t.Run(fmt.Sprintf("%v\n", test.name), func(t *testing.T) {
			f := filepath.Join("testdata/simplify", fmt.Sprintf("%v.igc", test.name))
			golden := fmt.Sprintf("%v.golden", f)
			track, err := ParseLocation(f)
			if err != nil {
				t.Fatal(err)
			}
			simple, err := track.Simplify(0.0001)
			if err != nil {
				t.Fatal(err)
			}

			// update golden if flag is passed
			if *update {
				jsn, _ := json.Marshal(simple)
				if err = ioutil.WriteFile(golden, jsn, 0644); err != nil {
					t.Fatal(err)
				}
			}

			b, _ := ioutil.ReadFile(golden)
			var goldenTrack Track
			json.Unmarshal(b, &goldenTrack)
			if len(simple.Points) != len(goldenTrack.Points) {
				t.Errorf("expected %v got %v simple points", len(goldenTrack.Points), len(simple.Points))
			}
		})
	}
}

func TestSimplifyStats(t *testing.T) {

	simplifyDB := "testdata/simplify/db"
	csvfile, _ := os.Create(fmt.Sprintf("%v/simplify-stats.csv", testDir))

	files, _ := ioutil.ReadDir(simplifyDB)
	fmt.Fprintf(csvfile, "id,total,clean,clean%%,simple001,simple001%%,simple0001,simple0001%%\n")
	for _, f := range files {
		track, _ := ParseLocation(fmt.Sprintf("%v/%v", simplifyDB, f.Name()))
		clean, err := track.Cleanup()
		if err != nil {
			t.Fatal(err)
		}
		if len(track.Points) == 0 {
			//t.Fatalf("track has no points :: %v", track)
		} else {
			// simplify using two different precisions
			simple001, err := clean.Simplify(0.001)
			if err != nil {
				t.Fatal(err)
			}
			simple0001, err := clean.Simplify(0.0001)
			if err != nil {
				t.Fatal(err)
			}
			// convert the whole lot to json for js usage
			jsn, _ := toLatLng(track)
			jsnClean, _ := toLatLng(clean)
			jsnSimple001, _ := toLatLng(simple001)
			jsnSimple0001, _ := toLatLng(simple0001)
			// generate the html/js content for visualization
			df := fmt.Sprintf("%v/js/%v.js", testDir, f.Name())
			d, _ := os.Create(df)
			d.WriteString(fmt.Sprintf("%s\n%s\n%s\n%s\n", append([]byte("path ="), jsn...), append([]byte("clean ="), jsnClean...), append([]byte("simple001 ="), jsnSimple001...), append([]byte("simple0001 ="), jsnSimple0001...)))
			d.Close()
			html := fmt.Sprintf(template, f.Name(), f.Name(), f.Name())
			ioutil.WriteFile(fmt.Sprintf("%v/%v.html", testDir, f.Name()), []byte(html), 0644)
			ptsClean := float64(len(clean.Points))
			// optimize for both simplified tracks
			opt := NewBruteForceOptimizer(false)
			//task001, err := opt.Optimize(simple001, 3, Distance)
			task001, err := Task{}, nil
			task0001, err := opt.Optimize(simple0001, 2, Distance)
			if err != nil {
				t.Fatal(err)
			}
			// fill in the csv line for this flight with simplify stats
			fmt.Fprintf(csvfile, "%v,%v,%v,%.1f,%v,%.1f,%v,%v,%.1f,%.1f\n",
				f.Name(),
				len(track.Points),
				len(clean.Points),
				float64(len(clean.Points))/ptsClean*100.0,
				len(simple001.Points),
				float64(len(simple001.Points))/ptsClean*100.0,
				task001.Distance(),
				len(simple0001.Points),
				float64(len(simple0001.Points))/ptsClean*100.0,
				task0001.Distance())
		}
	}
	csvfile.Close()
}

func toLatLng(track Track) ([]byte, error) {

	points := make([]s2.LatLng, len(track.Points))
	for i, v := range track.Points {
		points[i] = v.LatLng
	}
	return json.Marshal(points)
}

var template = `
<!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no">
    <meta charset="utf-8">
    <title>EzGliding Simplify</title>
    <style>
      /* Always set the map height explicitly to define the size of the div
       * element that contains the map. */
      #map {
        height: 100%%;
      }
      /* Optional: Makes the sample page fill the window. */
      html, body {
        height: 100%%;
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <div id="map"></div>
    <script
         src="file:///tmp/ezgliding/tests/js/%v.js"></script>

	<script>
      function initMap() {
        var map = new google.maps.Map(document.getElementById('map'), {
          zoom: 3,
          center: {lat: 0, lng: 0},
          mapTypeId: 'terrain'
        });

        for (i=0; i<path.length; i++) {
            path[i].lat = (path[i].Lat * 180) / Math.PI;
            path[i].lng = (path[i].Lng * 180) / Math.PI;
        }   
        for (i=0; i<clean.length; i++) {
            clean[i].lat = (clean[i].Lat * 180) / Math.PI;
            clean[i].lng = (clean[i].Lng * 180) / Math.PI;
        }   
        for (i=0; i<simple0001.length; i++) {
            simple0001[i].lat = (simple0001[i].Lat * 180) / Math.PI;
            simple0001[i].lng = (simple0001[i].Lng * 180) / Math.PI;
        }   
        for (i=0; i<simple001.length; i++) {
            simple001[i].lat = (simple001[i].Lat * 180) / Math.PI;
            simple001[i].lng = (simple001[i].Lng * 180) / Math.PI;
        }   
        console.log(path)
        var flightPath = new google.maps.Polyline({
          path: path,
          geodesic: true,
          strokeColor: '#FF0000',
          strokeOpacity: 0.5,
          strokeWeight: 2
        });
        var flightClean = new google.maps.Polyline({
          path: clean,
          geodesic: true,
          strokeColor: '#FFFF00',
          strokeOpacity: 0.5,
          strokeWeight: 2
        });
        var flightsimple0001 = new google.maps.Polyline({
          path: simple0001,
          geodesic: true,
          strokeColor: '#00FF00',
          strokeOpacity: 1.0,
          strokeWeight: 2
        });
        var flightsimple001 = new google.maps.Polyline({
          path: simple001,
          geodesic: true,
          strokeColor: '#0000FF',
          strokeOpacity: 1.0,
          strokeWeight: 2
        });

		flightClean.setMap(map);
        flightPath.setMap(map);
        flightsimple001.setMap(map);
        flightsimple0001.setMap(map);

		var bounds = new google.maps.LatLngBounds();
		var points = flightPath.getPath().getArray();
		for (var i = 0; i < points.length ; i++){
			bounds.extend(points[i]);
		}
		map.fitBounds(bounds);
      }
    </script>
    <script async defer
    src="https://maps.googleapis.com/maps/api/js?key=AIzaSyBFb6wMyglZopVA3DNX6gKM5gRYwWfwVAg&callback=initMap">
    </script>
  </body>
</html>
`
