#!/bin/bash

SCRIPT_PATH=$(dirname $(readlink -f $0))

function run_test()
{
  cmd=$1
  name=$2
  echo -e "\n *** $cmd $name ***"

  if ! $GOPATH/bin/grayt $cmd -o $SCRIPT_PATH/actual_$name.png; then
    echo -e "\nERROR"
    return 1
  fi

  cmp_output=$(compare -metric rmse $SCRIPT_PATH/actual_$name.png $SCRIPT_PATH/expect_$name.png null: 2>&1)
  echo output: $cmp_output
  if echo $cmp_output | grep "0 (0)"; then
    echo -e "\nPASSED"
    return 0
  else
    # Leave the diff behind to help debugging.
    compare $SCRIPT_PATH/actual_$name.png $SCRIPT_PATH/expect_$name.png $SCRIPT_PATH/diff_$name.png
    echo -e "\nFAILED"
    return 1
  fi
}

set -e
go install github.com/peterstace/grayt/cmd/grayt
run_test "-s cornellbox_classic     -d -w 1024 -q 1    " "classic_single"
run_test "-s cornellbox_classic     -d -w 128  -q 1000 " "classic"
run_test "-s spheretree             -d -w 256  -q 1    " "sphere_tree_single"
run_test "-s spheretree             -d -w 256  -q 100  " "sphere_tree"
run_test "-s splitbox               -d -w 128  -q 50   " "split_box"
run_test "-s cornellbox_reflections -d -w 64   -q 1000 " "reflection"
