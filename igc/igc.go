package igc

import (
    "fmt"
    "io/ioutil"
)

// Point is latitude/longitude (decimal format) and altitude (meters)
type Point struct {
    lat float32
    lon float32
    alt int
}

type Flight struct {
    bytes   []byte
    points  []Point
}

func (f *Flight) parse() {
    bytes, err := ioutil.ReadFile("test/basic.igc")
    f.bytes = bytes
}
