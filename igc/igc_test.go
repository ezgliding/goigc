package igc

import (
    "testing"
)

func TestBasic(t *testing.T) {
    flight, _ := Parse("flight-basic.igc")
    t.Logf("%s", flight)
}
