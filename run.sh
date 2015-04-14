#!/bin/bash

set -exu

ffmpeg -framerate 24 -i "out/%d.jpg" -codec copy "out.mkv"
