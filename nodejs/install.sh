#!/bin/bash

set -e

count=$(find . -type d -name "node_modules" | wc -l)

<<<<<<< 715e64a66d738b67d147cfd2087f234b031edbef
if [ "$count" -eq "0" ]; then
=======
if [ -d "$DIR" ]; then

	gulp dev
else
>>>>>>> update
	npm install
    gulp dev
else
	gulp dev
fi
