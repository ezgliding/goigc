#!/bin/bash
rc=0
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d | grep -v vendor); do
	if ls ${dir}/*.go &> /dev/null; then
		cd $dir
		go test -bench .
		if [ $? -ne 0 ]; then rc=1; fi
		cd -
	fi
done
exit $rc
