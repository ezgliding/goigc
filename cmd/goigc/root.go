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
	"os"

	"github.com/ezgliding/goigc/pkg/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "goigc",
		Short: "goigc is a parser and analyser for gliding flights",
		Long:  "",
		Version: fmt.Sprintf("%v %.7v %v %v", version.Version(), version.Commit(),
			version.BuildTime().Format("02/01/06 15:04:05"), version.Metadata()),
		Hidden:       true,
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceErrors, _ = cmd.Flags().GetBool("silent")
		},
	}
)

func init() {
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func main() {
	rootCmd.PersistentFlags().Bool("silent", false, "do not print any errors")
	rootCmd.SetVersionTemplate(
		`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}}
`)
	Execute()
}
