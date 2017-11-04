#!/usr/bin/env sh

set -ex

count1=$(find /etc/nginx/ssl/* -type d -name "intodevops.by" | wc -l)
count2=$(find /etc/nginx/conf.d/* -type f -name "ssl-intodevops$" | wc -l)

if [ "$count1" -gt "0" ] && [ "$count2" -gt "0" ]; then
	mv /etc/nginx/conf.d/ssl-intodevops /etc/nginx/conf.d/ssl-intodevops.conf
    sed -i 's/#return/return/g' /etc/nginx/conf.d/intodevops.conf
#    nginx -g "daemon off;"
else
:
#   nginx -g "daemon off;"
fi