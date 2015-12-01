#!/usr/bin/env bash

if [ -z "$1" ]; then
    IMAGE="ethernetdan/wedirect"
    echo "Ran without parameters, using default image name"
else
    IMAGE=$1
fi

echo "Building wedirect image as $IMAGE"
echo

docker run --rm \
    -v "$(pwd):/src" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    centurylink/golang-builder \
    $IMAGE
