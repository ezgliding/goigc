#!/bin/bash
set -e
set -x
go test -bench=Benchmark* -cpuprofile=cpu.pprof -memprofile=mem.pprof -run=x
go tool pprof -top -cum cpu.pprof
go tool pprof -top -cum mem.pprof
