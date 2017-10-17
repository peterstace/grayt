#!/bin/bash

set -e

go test -v ./...
./systemtest/run.sh
