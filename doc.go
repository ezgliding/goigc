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

// Package goigc contains functionality to parse and analyse gliding flights.
//
// goigc contains the following packages:
//
// The cmd packages contains command line utilities, namely the goigc binary.
//
// The pkg/igc provides the parsing and analysis code for the igc format.
//
package goigc

// blank imports help docs.
import (
	// pkg/igc package
	_ "github.com/ezgliding/goigc/pkg/igc"
	// pkg/version package
	_ "github.com/ezgliding/goigc/pkg/version"
)
