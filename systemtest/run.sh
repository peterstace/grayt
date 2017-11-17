#!/bin/bash

SCRIPT_PATH=$(dirname $(readlink -f $0))

compare --version

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
    compare $SCRIPT_PATH/actual_$name.png $SCRIPT_PATH/expect_$name.png $SCRIPT_PATH/diff_$name.png
    echo -e "\nFAILED"
    return 1
  fi
}

set -e
run_test "go run -race examples/cornellbox/classic/main.go       -d -w 1024 -q 1"    "classic_single"
run_test "go run       examples/cornellbox/classic/main.go       -d -w 128  -q 1000" "classic"
run_test "go run       examples/cornellbox/spheretree/main.go    -d -w 256  -q 1"    "sphere_tree_single"
run_test "go run       examples/cornellbox/spheretree/main.go    -d -w 256  -q 100"  "sphere_tree"
run_test "go run       examples/cornellbox/splitbox/main.go      -d -w 128  -q 50"   "split_box"
run_test "go run       examples/cornellbox/reflections/main.go   -d -w 64   -q 1000" "reflection"
