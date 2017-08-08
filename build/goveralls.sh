#!/bin/bash
golint .
go vet .
echo "mode: count" > profile.cov
go test -covermode=count -coverprofile=profile.cov .
goveralls -coverprofile=profile.cov -service=travis-ci -repotoken $COVERALLS_TOKEN
