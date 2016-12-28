#!/bin/bash

SCRIPT_PATH=$(dirname $(readlink -f $0))

function run_test()
{
  cmd=$1
  name=$2

  echo -e "\n *** $name ***"

  if $cmd -o $SCRIPT_PATH/actual_$name.png &&
    compare -metric rmse $SCRIPT_PATH/actual_$name.png $SCRIPT_PATH/expect_$name.png null: ; then
    echo -e "\nPASSED"
    return 0
  else
    echo -e "\nFAILED"
    return 1
  fi
}

set -e
run_test "go run examples/cornellbox/classic/main.go    -w 512 -q 1"    "classic_single"
run_test "go run examples/cornellbox/classic/main.go    -w 128 -q 1000" "classic"
run_test "go run examples/cornellbox/spheretree/main.go -w 256 -q 100"  "sphere_tree"
run_test "go run examples/cornellbox/splitbox/main.go   -w 128 -q 50"   "split_box"
