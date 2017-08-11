#!/bin/bash
set -e
set -x
export TARGET=${1:-master}
export PATH=/tmp/perf/cmd/benchstat:$PATH
if [ ! -f /tmp/perf/cmd/benchstat/benchstat ]; then
	go get golang.org/x/perf/benchstat
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
git config --replace-all remote.origin.fetch +refs/heads/*:refs/remotes/origin/*
git fetch
git fetch --tags
git checkout -f $TARGET &> /dev/null
go test -v -bench=Benchmark* -run None > bench-${TARGET}.result

benchstat -delta-test none $TARGET_RESULT $BRANCH_RESULT
rm -f $TARGET_RESULT $BRANCH_RESULT

git checkout -f $BRANCH &> /dev/null
