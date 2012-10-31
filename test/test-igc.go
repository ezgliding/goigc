package main

import (
    "fmt"
)

func main() {
    f := Flight()
    f.parse()
    fmt.Printf(f.bytes)
}
