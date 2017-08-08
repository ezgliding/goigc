#!/bin/bash
rc=0
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d | grep -v vendor); do
	if ls ${dir}/*.go &> /dev/null; then
		cd $dir
		go test -bench . -cpuprofile=cpu.pprof -memprofile=mem.pprof -run=Benchmark*
		if [ $? -ne 0 ]; then rc=1; fi
		if [ -f ${dir}.test ] && [ -f cpu.pprof ]; then
			go tool pprof -text ${dir}.test cpu.pprof
			if [ $? -ne 0 ]; then rc=1; fi
		fi
		if [ -f ${dir}.test ] && [ -f mem.pprof ]; then
			go tool pprof -text ${dir}.test mem.pprof
			if [ $? -ne 0 ]; then rc=1; fi
		fi
		rm -f ${dir}.test cpu.pprof mem.pprof
		cd -
	fi
done
exit $rc
