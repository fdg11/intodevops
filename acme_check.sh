#!/usr/bin/env bash

set -e

count=$(find nginx/ssl/* -type d -name "intodevops.by" | wc -l)

if [ "$count" -eq "0" ]; then
    docker-compose exec proxy acme.sh --issue -d intodevops.by -d www.intodevops.by -k 4096 -w /workspace
cat <<EOF > nginx/conf.d/ssl-intodevops.conf
    server {
    listen 443 ssl http2;
    server_name intodevops.by www.intodevops.by;

    return 301 \$scheme://intodevops.by\$request_uri;
    location /.well-known {
    root /workspace;
    }

    ssl_certificate ssl/intodevops.by/fullchain.cer;
    ssl_certificate_key ssl/intodevops.by/intodevops.by.key;
    include ssl/ssl.conf;

    location / {
        proxy_pass http://go:8080;
        }
    }
EOF
    sed -i 's/#return/return/g' nginx/conf.d/intodevops.conf
    docker-compose exec proxy nginx -s reload
else
    docker-compose exec proxy nginx -s reload
fi