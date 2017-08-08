#!/bin/bash
go test -bench . -cpuprofile=cpu.pprof -memprofile=mem.pprof -run=Benchmark*
go tool pprof -text cpu.pprof
go tool pprof -text mem.pprof
