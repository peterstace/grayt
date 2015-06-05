#!/bin/bash

set -ex

dt=`date "+%Y-%m-%d_%H%M%S"`
rev=`git rev-parse --short HEAD`

mkdir -p img

$GOPATH/bin/main -o img/$dt\_$rev\_10.png  -s 10
$GOPATH/bin/main -o img/$dt\_$rev\_100.png -s 100
$GOPATH/bin/main -o img/$dt\_$rev\_1000.png -s 1000
