#!/bin/bash

set -ex

dt=`date "+%Y-%m-%d_%H%M%S"`
rev=`git rev-parse --short HEAD`

mkdir -p img

$GOPATH/bin/main -o img/$dt\_$rev\_16.png   -s 16
$GOPATH/bin/main -o img/$dt\_$rev\_32.png   -s 32
$GOPATH/bin/main -o img/$dt\_$rev\_64.png   -s 64
$GOPATH/bin/main -o img/$dt\_$rev\_128.png  -s 128
$GOPATH/bin/main -o img/$dt\_$rev\_256.png  -s 256
$GOPATH/bin/main -o img/$dt\_$rev\_512.png  -s 512
$GOPATH/bin/main -o img/$dt\_$rev\_1024.png -s 1024
$GOPATH/bin/main -o img/$dt\_$rev\_2048.png -s 2048
$GOPATH/bin/main -o img/$dt\_$rev\_4096.png -s 4096
$GOPATH/bin/main -o img/$dt\_$rev\_8192.png -s 8192
