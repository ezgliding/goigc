#!/bin/bash
go test -bench=Benchmark* -cpuprofile=cpu.pprof -memprofile=mem.pprof -run=x
go tool pprof -text cpu.pprof
go tool pprof -text mem.pprof
