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
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ezgliding/goigc/pkg/igc"
	"github.com/ezgliding/goigc/pkg/netcoupe"
)

func init() {
	crawlCmd.Flags().String("source", "netcoupe", "online web source to crawl")
	rootCmd.AddCommand(crawlCmd)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl START END PATH",
	Short: "crawls flights from the given web source",
	Long: `Crawls the given web source for gliding flights between START and END date.

Expected format for start and end dates is 2006-01-02.

Results are stored under PATH with the following structure.

PATH/YEAR
  /DD-MM-YYYY.json ( one json file per day with flight metadata )
  /flights
    /TRACKID ( one file with the flight track in the original format )
`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {

		start, err := time.Parse("2006-01-02", args[0])
		if err != nil {
			return err
		}
		end, err := time.Parse("2006-01-02", args[1])
		if err != nil {
			return err
		}
		if start.Year() != end.Year() {
			return fmt.Errorf("Start and end year must be the same")
		}
		basePath := args[2]
		err = os.MkdirAll(fmt.Sprintf("%v/%v/flights", basePath, start.Year()), os.ModePerm)
		if err != nil {
			return err
		}

		var n netcoupe.Netcoupe = netcoupe.NewNetcoupeYear(start.Year())
		current := start
		for ; end.After(current.AddDate(0, 0, -1)); current = current.AddDate(0, 0, 1) {
			var flights []igc.Flight
			dbFile := fmt.Sprintf("%v/%v/%v.json", basePath, current.Year(), current.Format("02-01-2006"))
			if _, err := os.Stat(dbFile); os.IsNotExist(err) {
				flights, err = n.Crawl(current, current)
				if err != nil {
					return err
				}
				jsonFlights, err := json.MarshalIndent(flights, "", "   ")
				if err != nil {
					return err
				}
				err = ioutil.WriteFile(dbFile, jsonFlights, 0644)
				if err != nil {
					return err
				}
			} else {
				b, err := ioutil.ReadFile(dbFile)
				if err != nil {
					return err
				}
				err = json.Unmarshal(b, &flights)
				if err != nil {
					return err
				}
			}

			for _, f := range flights {
				flightFile := fmt.Sprintf("%v/%v/flights/%v", basePath, current.Year(), f.TrackID)
				if _, err := os.Stat(flightFile); os.IsNotExist(err) {
					url := fmt.Sprintf("%v%v", n.TrackBaseUrl(), f.TrackID)
					data, err := n.Get(url)
					if err != nil {
						return err
					}
					err = ioutil.WriteFile(flightFile, data, 0644)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	},
}
