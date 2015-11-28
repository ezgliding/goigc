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

	"github.com/spf13/cobra"
)

var version = "0.1"

// versionCmd respresents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of goigc",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("version: %v\n", version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings

	// Cobra supports Persistent Flags which will work for this command and all subcommands
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle" )

}
