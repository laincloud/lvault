#!/bin/bash

mkdir -p /go/src/github.com/laincloud/
ln -sf /lain/app /go/src/github.com/laincloud/lvault
gobuildweb dist
unzip -o /lain/app/lvault-0.1.zip -d /lain/app
cd /lain/app/lvault-0.1

