#!/usr/bin/env sh

set -ex

count=$(find /etc/nginx/ssl/* -type d -name "intodevops.by" | wc -l)

if [ "$count" -gt "0" ]; then
	mv /etc/nginx/conf.d/ssl-intodevops /etc/nginx/conf.d/ssl-intodevops.conf
    sed -i 's/#return/return/g' /etc/nginx/conf.d/intodevops.conf
fi