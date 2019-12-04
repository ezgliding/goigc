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
	// cmd/goigc package
	_ "github.com/ezgliding/goigc/cmd/goigc"
	// pkg/igc package
	_ "github.com/ezgliding/goigc/pkg/igc"
	// pkg/version package
	_ "github.com/ezgliding/goigc/pkg/version"
)
