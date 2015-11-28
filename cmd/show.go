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

package cmd

import (
	"log"
	"os"
	"text/template"

	"github.com/rochaporto/goigc/igc"
	"github.com/spf13/cobra"
)

var (
	all bool
)

// showCmd respresents the show command
var showCmd = &cobra.Command{
	Use:   "show [path]",
	Short: "Show track details",
	Long:  `Parse track at a given path and show header and stats data.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("no track location given")
		}
		t, err := igc.ParseLocation(args[0])
		if err != nil {
			log.Fatalf("parse failed: %v\n", err)
		}
		tmpl, err := template.New("show").Parse(`
Pilot: {{.Pilot}}
`)
		if err != nil {
			log.Fatalf("template parse failed: %v\n", err)
		}
		err = tmpl.Execute(os.Stdout, t)
		if err != nil {
			log.Fatalf("show display failed: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(showCmd)

	// Command flags
	showCmd.PersistentFlags().BoolVarP(&all, "all", "a", true, "Show all (including calculated) flight data")
}
