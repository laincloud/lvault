#!/bin/bash

go get github.com/tools/godep
mkdir -p $GOPATH/src/github.com/hashicorp
cd $GOPATH/src/github.com/hashicorp
git clone https://github.com/hashicorp/vault.git
# go get github.com/hashicorp/vault
cd vault
git checkout v0.9.5
make bootstrap
make
make dev
go get github.com/mijia/sweb/log

cd /lain/app

yum -y install unzip
git clone --depth=1 https://github.com/golang/sys.git /go/src/golang.org/x/sys # GFW
git clone --depth=1 https://github.com/golang/net.git /go/src/golang.org/x/net
go get github.com/mijia/gobuildweb
cd $GOPATH/src/github.com/mijia/gobuildweb && sed -i '/deps = append(deps, "browserify", "coffeeify", "envify", "uglifyify", "babelify", "babel-preset-es2015", "babel-preset-react", "nib", "stylus")/d' cmds.go  && go install
gobuildweb dist
ls -1 | grep -v node_modules | xargs rm -rf
