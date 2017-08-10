#!/bin/bash
export TARGET=${1:-master}
export PATH=/tmp/perf/cmd/benchstat:$PATH
if [ ! -f /tmp/perf/cmd/benchstat/benchstat ]; then
	rm -rf /tmp/perf
	git clone https://github.com/golang/perf.git /tmp/perf
	cd /tmp/perf/cmd/benchstat
	go build .
	cd -
fi

REVBRANCH=$(git rev-parse --abbrev-ref HEAD)
BRANCH=${TRAVIS_BRANCH:-$REVBRANCH}
TARGET_RESULT="bench-${TARGET}.result"
BRANCH_RESULT="bench-${BRANCH}.result"
go test -bench=Benchmark* -run None > bench-${BRANCH}.result
git checkout $TARGET &> /dev/null
go test -v -bench=Benchmark* -run None > bench-${TARGET}.result

benchstat -delta-test none $TARGET_RESULT $BRANCH_RESULT
rm $TARGET_RESULT $BRANCH_RESULT

git checkout $BRANCH &> /dev/null
