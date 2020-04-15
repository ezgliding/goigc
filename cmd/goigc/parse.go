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
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/ezgliding/goigc/pkg/igc"
)

func init() {
	// TODO(rochaporto): not yet supported, only igc
	parseCmd.Flags().String("format", "", "input file format - auto detection by default")
	parseCmd.Flags().Bool("no-points", false, "do not include individual points")
	parseCmd.Flags().String("output-format", "yaml", "output format for display")
	parseCmd.Flags().String("output-file", "/dev/stdout", "output file to write to")
	rootCmd.AddCommand(parseCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse FILE",
	Short: "parses information about the given flight",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFile, err := cmd.Flags().GetString("output-file")
		if err != nil {
			return err
		}
		outputFormat, err := cmd.Flags().GetString("output-format")
		if err != nil {
			return err
		}

		trk, err := igc.ParseLocation(args[0])
		if err != nil {
			return err
		}

		noPoints, _ := cmd.Flags().GetBool("no-points")
		if noPoints {
			trk = igc.Track{Header: trk.Header}
		}
		result, err := trk.Encode(outputFormat)
		if err != nil {
			return err
		}
		if outputFile == "/dev/stdout" {
			fmt.Printf("%v", string(result))
		} else {
			err = ioutil.WriteFile(outputFile, result, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
