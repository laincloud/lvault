#!/usr/bin/env bash

set -e
export GOPATH=/go
cd $GOPATH/src/github.com/laincloud/lvault
echo "" > coverage.txt

for d in $(godep go list ./...); do
    godep go test -v -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
