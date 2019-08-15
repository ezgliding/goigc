// Copyright 2017 The ezgliding Authors.
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

/*
Package igc provides means to parse and analyse files in the IGC format.

This format is defined by the International Gliding Commission (IGC) and
was created to set a standard for recording gliding flights.

The full specification is available in Appendix A of the IGC FR Specification:
http://www.fai.org/component/phocadownload/category/?download=11005

Calculation of the optimal flight distance considering multiple turnpoints and
FAI triangles are available via Optimizers. Available Optimizers include brute
force, montecarlo method, genetic algorithms, etc.

*/
package igc
