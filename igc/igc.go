package igc

import (
    "bufio"
    "io"
    "os"
)

// Point is latitude/longitude (decimal format) and altitude (meters)
type Point struct {
    lat float32
    lon float32
    alt int
}

// Track is a list of Points representing a Flight track
type Track []Point

// Flight has all the flight info
type Flight struct {
    bytes   []byte
    track   Track
    record  []byte
}

// Parses the IGC file at the given location, returning a Flight
func Parse(location string) (flight Flight, err error) {
    fd, _ := os.Open(location)
    defer fd.Close()
    r := bufio.NewReader(fd)
    for record, err := r.ReadSlice('\n'); err != io.EOF; record, err = r.ReadSlice('\n') {

        switch record[0] {
        case 'A':
        case 'B':
        case 'C':
        }
    }
    return
}

