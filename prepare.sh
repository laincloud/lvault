#!/bin/bash

go get github.com/tools/godep
mkdir -p $GOPATH/src/github.com/hashicorp
cd $GOPATH/src/github.com/hashicorp
git clone https://github.com/hashicorp/vault.git
cd vault
make bootstrap
make
make dev
go get github.com/mijia/sweb/log

cd /lain/app

yum -y install unzip
git clone --depth=1 https://github.com/golang/sys.git /go/src/golang.org/x/sys # GFW
git clone --depth=1 https://github.com/golang/net.git /go/src/golang.org/x/net
go get github.com/mijia/gobuildweb
gobuildweb dist
ls -1 | grep -v node_modules | xargs rm -rf
