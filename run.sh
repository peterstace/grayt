#!/bin/bash

set -e

uuid=$(curl -s localhost:6060/renders -XPOST | jq -r .uuid)

showStatus() {
	json=$(curl -s localhost:6060/renders/$1 -XGET)
	echo $json | jq .
}

showStatus $uuid

curl -s localhost:6060/renders/$uuid/scene -d 'cornellbox_classic' -XPUT
curl -s localhost:6060/renders/$uuid/running -d 'true' -XPUT
