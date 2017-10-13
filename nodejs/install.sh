#!/bin/bash

set -e

count=$(find /workspace -type d -name "node_modules" | wc -l)

if [ "$count" eq 0 ]; then
	npm install
    gulp dev
else
	gulp dev
fi
