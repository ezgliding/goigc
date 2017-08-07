#!/bin/bash
rc=0
export PATH=${PATH}:${HOME}/gopath/bin
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d); do
	if ls ${dir}/*.go &> /dev/null; then
		cd $dir
		go test -bench . -cpuprofile=cpu.pprof -memprofile=mem.pprof -run=Benchmark*
		if [ -f ${dir}.test ] && [ -f cpu.pprof ]; then
			go tool pprof -text ${dir}.test cpu.pprof | head -n 15
		fi
		if [ -f ${dir}.test ] && [ -f mem.pprof ]; then
			go tool pprof -text ${dir}.test mem.pprof | head -n 15
		fi
		rm -f ${dir}.test cpu.pprof mem.pprof
		cd -
	fi
done
if [ $? -ne 0 ]; then rc=1; fi
exit $rc
