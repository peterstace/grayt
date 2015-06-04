#!/bin/bash

set -ex

dt=`date "+%Y-%m-%d_%H%M%S"`
rev=`git rev-parse --short HEAD`

mkdir -p img

for i in `seq 1 3`;
do
	k=$((10 ** i))
	$GOPATH/bin/main -o img/$dt\_$rev\_$k.png -s $k
done
