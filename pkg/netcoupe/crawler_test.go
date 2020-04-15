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

package netcoupe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ezgliding/goigc/pkg/igc"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func TestNetcoupeCrawler(t *testing.T) {
	start := time.Date(2018, time.December, 24, 12, 0, 0, 0, time.UTC)
	end := time.Date(2018, time.December, 25, 12, 0, 0, 0, time.UTC)

	var n Netcoupe = NewNetcoupeYear(2018)
	flights, err := n.Crawl(start, end)
	if err != nil {
		t.Errorf("%v", err)
	}

	if len(flights) <= 0 {
		t.Errorf("no flights returned")
	}

	jsonFlights, _ := json.MarshalIndent(flights, "", "   ")
	fmt.Printf("%v\n", string(jsonFlights))
}

func TestNetcoupeCrawlerDownload(t *testing.T) {

	year := 2019
	start := time.Date(year, 12, 31, 12, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 12, 0, 0, 0, time.UTC)

	var n Netcoupe = NewNetcoupeYear(year)
	current := start
	for ; end.After(current.AddDate(0, 0, -1)); current = current.AddDate(0, 0, 1) {
		var flights []igc.Flight
		dbFile := fmt.Sprintf("db/%v/%v.json", year, current.Format("02-01-2006"))
		if _, err := os.Stat(dbFile); os.IsNotExist(err) {
			flights, err = n.Crawl(current, current)
			if err != nil {
				t.Errorf("%v", err)
			}
			jsonFlights, _ := json.MarshalIndent(flights, "", "   ")
			// TODO(rochaporto): handle error
			_ = ioutil.WriteFile(dbFile, jsonFlights, 0644)
		} else {
			b, _ := ioutil.ReadFile(dbFile)
			err = json.Unmarshal(b, &flights)
			if err != nil {
				fmt.Printf("error parsing %v %v\n", dbFile, err)
				continue
			}
		}

		for _, f := range flights {
			flightFile := fmt.Sprintf("db/%v/flights/%v", year, f.TrackID)
			if _, err := os.Stat(flightFile); os.IsNotExist(err) {
				url := fmt.Sprintf("%v%v", n.TrackBaseUrl(), f.TrackID)
				data, _ := n.Get(url)
				// TODO(rochaporto): handle error
				_ = ioutil.WriteFile(flightFile, data, 0644)
			}
		}
	}
}
