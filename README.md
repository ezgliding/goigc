# goigc [![Build Status](https://travis-ci.org/ezgliding/goigc.svg?branch=master)](http://travis-ci.org/ezgliding/goigc) [![Coverage Status](https://coveralls.io/repos/github/ezgliding/goigc/badge.svg?branch=master)](https://coveralls.io/github/ezgliding/goigc?branch=vendor) [![GoDoc](https://godoc.org/github.com/ezgliding/goigc?status.png)](https://godoc.org/github.com/ezgliding/goigc) [![Go Report Card](https://goreportcard.com/badge/github.com/ezgliding/goigc)](https://goreportcard.com/report/github.com/ezgliding/goigc) ![Project Status](http://img.shields.io/badge/status-prealpha-red.svg)

Handles flight tracks in [IGC](http://www.fai.org/component/phocadownload/category/?download=5745:igc-flight-recorder-specification-edition-2-with-al1-2011-5-31) format.

## Usage

    $ go get github.com/rochaporto/goigc

    $ ./goigc 
    $ Parse and analyse flight tracks in IGC format.
    $                          _______
    $                             |
    $ /--------------------------(_)--------------------------\
    $ 
    $ Usage:
    $   goigc [command]
    $ 
    $ Available Commands:
    $   convert     Convert track between different formats
    $   optimize    Optimize the track (for distance and score)
    $   show        Show track details
    $   version     Print version of goigc
    $ 
    $ Flags:
    $       --config="": config file (default is $HOME/.goigc.yaml)
    $   -h, --help[=false]: help for goigc
    $ 
    $ Use "goigc [command] --help" for more information about a command.

    $ goigc stats sample-flight.igc

    $ goigc optimize sample-flight.igc

## Testing

Tests rely on the golden test pattern. To update the test data under the trest
directory run the tests with the `update` flag:
```
go test -update .
```

## Documentation

    $ godoc github.com/rochaporto/goigc
