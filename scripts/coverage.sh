#!/usr/bin/env bash
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
# Based on the Helm script file:
# github.com/helm/helm/scripts/coverage.sh

set -euo pipefail

covermode=${COVERMODE:-atomic}
coverdir=$(mktemp -d /tmp/coverage.XXXXXXXXXX)
profile="${coverdir}/cover.out"

pushd /
hash goveralls 2>/dev/null || go get github.com/mattn/goveralls
popd

generate_cover_data() {
  for d in $(go list ./...) ; do
    (
      local output="${coverdir}/${d//\//-}.cover"
      go test -coverprofile="${output}" -covermode="$covermode" "$d"
    )
  done

  echo "mode: $covermode" >"$profile"
  grep -h -v "^mode:" "$coverdir"/*.cover >>"$profile"
}

push_to_coveralls() {
  goveralls -coverprofile="${profile}" -service=circle-ci
}

generate_cover_data
go tool cover -func "${profile}"

case "${1-}" in
  --html)
    go tool cover -html "${profile}"
    ;;
  --coveralls)
    push_to_coveralls
    ;;
esac

