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
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ezgliding/goigc/pkg/igc"
	"github.com/spf13/cobra"
)

func init() {
	phasesCmd.Flags().String("output-format", "yaml", "output format for display")
	phasesCmd.Flags().String("output-file", "", "output file to write to")
	rootCmd.AddCommand(phasesCmd)
}

var phasesCmd = &cobra.Command{
	Use:   "phases FILE",
	Short: "compute phases for the given flight",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		/**	outputFile, err := cmd.Flags().GetString("output-file")
		if err != nil {
			return err
		}
		outputFormat, err := cmd.Flags().GetString("output-format")
		if err != nil {
			return err
		}
		*/
		trk, err := igc.ParseLocation(args[0])
		if err != nil {
			return err
		}
		phases, err := trk.Phases()
		if err != nil {
			return err
		}
		records := make([][]string, len(phases)+1)
		records[0] = []string{
			"Flight", "Type", "CirclingType", "StartTime", "StartAlt",
			"StartIndex", "EndTime", "EndAlt", "EndIndex", "Duration",
			"AvgVario", "TopVario", "AvgGndSpeed", "TopGndSpeed", "Distance",
			"LD", "CentroidLat", "CentroidLng", "CellID"}
		var p igc.Phase
		for i := 0; i < len(phases); i++ {
			p = phases[i]
			records[i+1] = []string{
				"",
				fmt.Sprintf("%d", p.Type),
				fmt.Sprintf("%d", p.CirclingType),
				fmt.Sprintf("%v", p.Start.Time.Format("15:04:05")),
				fmt.Sprintf("%d", p.Start.GNSSAltitude),
				fmt.Sprintf("%d", p.StartIndex),
				fmt.Sprintf("%v", p.End.Time.Format("15:04:05")),
				fmt.Sprintf("%d", p.End.GNSSAltitude),
				fmt.Sprintf("%d", p.EndIndex),
				fmt.Sprintf("%f", p.Duration().Seconds()),
				fmt.Sprintf("%f", p.AvgVario), fmt.Sprintf("%f", p.TopVario),
				fmt.Sprintf("%f", p.AvgGndSpeed),
				fmt.Sprintf("%f", p.TopGndSpeed),
				fmt.Sprintf("%f", p.Distance), fmt.Sprintf("%f", p.LD),
				fmt.Sprintf("%f", p.Centroid.Lat.Degrees()),
				fmt.Sprintf("%f", p.Centroid.Lng.Degrees()),
				fmt.Sprintf("%d", p.CellID)}
		}

		w := csv.NewWriter(os.Stdout)
		err = w.WriteAll(records)
		if err != nil {
			return err
		}

		return nil
	},
}
