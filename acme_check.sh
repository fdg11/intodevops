#!/usr/bin/env bash

set -e

count=$(find nginx/ssl/* -type d -name "intodevops.by" | wc -l)

if [ "$count" -eq "0" ]; then
    docker-compose exec proxy acme.sh --issue -d intodevops.by -d www.intodevops.by -d sub.intudevops.by -k 4096 -w /workspace --force
cat <<EOF > nginx/conf.d/ssl-intodevops.conf
    upstream portainer {
        server portainer:9000;
    }

    server {
        listen 443 ssl http2;
        server_name www.intodevops.by;
        return 301 https://intodevops.by\$request_uri;

        ssl_certificate ssl/intodevops.by/fullchain.cer;
        ssl_certificate_key ssl/intodevops.by/intodevops.by.key;
        include ssl/ssl.conf;

    }

    server {
        listen 443 ssl http2;
        server_name intodevops.by;

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

    server {
        listen 443 ssl http2;
        server_name sub.intodevops.by;

        ssl_certificate ssl/intodevops.by/fullchain.cer;
        ssl_certificate_key ssl/intodevops.by/intodevops.by.key;
        include ssl/ssl.conf;

        location /portainer/ {
            proxy_http_version 1.1;
            proxy_set_header Host              \$http_host;   # required for docker client's sake
            proxy_set_header X-Real-IP         \$remote_addr; # pass on real client's IP
            proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
            proxy_read_timeout                 900;

            proxy_set_header Connection "";
            proxy_buffers 32 4k;
            proxy_pass http://portainer/;
        }

        location /portainer/api/websocket/ {
            proxy_http_version 1.1;
            proxy_set_header Upgrade \$http_upgrade;
            proxy_set_header Connection \$connection_upgrade;
            proxy_pass http://portainer/api/websocket/;
        }

    }
EOF
    sed -i 's/#return/return/g' nginx/conf.d/intodevops.conf
    docker-compose exec proxy nginx -s reload
else
    docker-compose exec proxy nginx -s reload
fi