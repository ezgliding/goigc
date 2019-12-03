#!/bin/bash
#
# Copyright The ezgliding Authors.
# 
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# 
set -e
set -x
export TARGET=${1:-master}
export PATH=/tmp/perf/cmd/benchstat:$PATH
export GO111MODULE=on
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
git config --replace-all remote.origin.fetch +refs/heads/*:refs/remotes/origin/*
git fetch
git fetch --tags
git checkout -f $TARGET &> /dev/null
go test -v -bench=Benchmark* -run None > bench-${TARGET}.result

benchstat -delta-test none $TARGET_RESULT $BRANCH_RESULT
rm -f $TARGET_RESULT $BRANCH_RESULT

git checkout -f $BRANCH &> /dev/null
