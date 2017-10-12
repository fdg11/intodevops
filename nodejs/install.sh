#!/bin/bash

set -e

DIR=/workspace/node-modules

if [ -d "$DIR" ]; then
	gulp dev
else
	npm install
	gulp dev
fi
