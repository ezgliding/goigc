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
	"path/filepath"
	"testing"

	"github.com/golang/geo/s2"
)

type simplifyTest struct {
	name   string
	result string
}

var simplifyTests = []simplifyTest{
	{
		name:   "simplify-short-flight-1",
		result: "simplify-short-flight-1.simplified",
	},
}

func TestSimplify(t *testing.T) {
	for _, test := range simplifyTests {
		t.Run(fmt.Sprintf("%v\n", test.name), func(t *testing.T) {
			f := filepath.Join("testdata", fmt.Sprintf("%v.igc", test.name))
			golden := fmt.Sprintf("%v.golden", f)
			track, err := ParseLocation(f)
			if err != nil {
				t.Fatal(err)
			}
			simplified := track.Simplify(0.0001)

			// update golden if flag is passed
			if *update {
				jsn, _ := json.Marshal(simplified)
				if err = ioutil.WriteFile(golden, jsn, 0644); err != nil {
					t.Fatal(err)
				}
			}

			b, _ := ioutil.ReadFile(golden)
			var goldenTrack Track
			json.Unmarshal(b, &goldenTrack)
			if len(simplified.Points) != len(goldenTrack.Points) {
				t.Errorf("expected %v got %v simplified points", len(goldenTrack.Points), len(simplified.Points))
			}
		})
	}
}

func TestSimplifyCompare(t *testing.T) {
	dir := "/home/ricardo/ws/ezgliding/crawler/db/2018/flights"
	dest := "/home/ricardo/tests"
	files, _ := ioutil.ReadDir(dir)
	fmt.Printf("id,original,cleanup,simplified001,simplified0001")
	for _, f := range files {
		track, _ := ParseLocation(fmt.Sprintf("%v/%v", dir, f.Name()))
		clean := track.Cleanup()
		if len(track.Points) == 0 {
		} else {
			jsn, _ := toLatLng(track)
			simplified001 := clean.Simplify(0.001)
			simplified0001 := clean.Simplify(0.0001)
			jsnSimplified001, _ := toLatLng(simplified001)
			jsnSimplified0001, _ := toLatLng(simplified0001)
			ioutil.WriteFile(fmt.Sprintf("%v/js/%v.js", dest, f.Name()), append([]byte("path ="), jsn...), 0644)
			ioutil.WriteFile(fmt.Sprintf("%v/js/%v-simplified001.js", dest, f.Name()), append([]byte("simplified001 = "), jsnSimplified001...), 0644)
			ioutil.WriteFile(fmt.Sprintf("%v/js/%v-simplified0001.js", dest, f.Name()), append([]byte("simplified0001 = "), jsnSimplified0001...), 0644)

			html := fmt.Sprintf(template, f.Name(), f.Name(), f.Name())
			ioutil.WriteFile(fmt.Sprintf("%v/%v.html", dest, f.Name()), []byte(html), 0644)
			fmt.Printf("%v,%v,%v,%v,%v\n", f.Name(), len(track.Points), len(clean.Points), len(simplified001.Points), len(simplified0001.Points))
		}
	}
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
    <title>Simple Polylines</title>
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
         src="file:///home/ricardo/tests/js/%v.js"></script>
    <script
         src="file:///home/ricardo/tests/js/%v-simplified001.js"></script>
    <script
         src="file:///home/ricardo/tests/js/%v-simplified0001.js"></script>
    <script>

      // This example creates a 2-pixel-wide red polyline showing the path of
      // the first trans-Pacific flight between Oakland, CA, and Brisbane,
      // Australia which was made by Charles Kingsford Smith.

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
        for (i=0; i<simplified0001.length; i++) {
            simplified0001[i].lat = (simplified0001[i].Lat * 180) / Math.PI;
            simplified0001[i].lng = (simplified0001[i].Lng * 180) / Math.PI;
        }   
        for (i=0; i<simplified001.length; i++) {
            simplified001[i].lat = (simplified001[i].Lat * 180) / Math.PI;
            simplified001[i].lng = (simplified001[i].Lng * 180) / Math.PI;
        }   
        console.log(path)
        var flightPath = new google.maps.Polyline({
          path: path,
          geodesic: true,
          strokeColor: '#FF0000',
          strokeOpacity: 1.0,
          strokeWeight: 2
        });
        var flightSimplified0001 = new google.maps.Polyline({
          path: simplified0001,
          geodesic: true,
          strokeColor: '#00FF00',
          strokeOpacity: 1.0,
          strokeWeight: 2
        });
        var flightSimplified001 = new google.maps.Polyline({
          path: simplified001,
          geodesic: true,
          strokeColor: '#0000FF',
          strokeOpacity: 1.0,
          strokeWeight: 2
        });

        flightPath.setMap(map);
        flightSimplified0001.setMap(map);
        flightSimplified001.setMap(map);

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
