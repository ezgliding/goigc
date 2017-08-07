#!/bin/bash
rc=0
export PATH=$PATH:$HOME/gopath/bin
golint .
if [ $? -ne 0 ]; then rc=1; fi
go vet ./...
if [ $? -ne 0 ]; then rc=1; fi
echo "mode: count" > profile.cov
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d); do
	if ls $dir/*.go &> /dev/null; then
		go test -covermode=count -coverprofile=$dir/profile.tmp $dir
		if [ $? -ne 0 ]; then rc=1; fi
		if [ -f $dir/profile.tmp ]; then
			cat $dir/profile.tmp | tail -n +2 >> profile.cov
			rm $dir/profile.tmp
		fi
	fi
done
goveralls -coverprofile=profile.cov -service=travis-ci -repotoken $COVERALLS_TOKEN
if [ $? -ne 0 ]; then rc=1; fi
exit $rc
