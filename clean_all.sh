#!/usr/bin/env bash

echo -e "Clean All data docker!"
docker system prune -f
imgs=$(docker images -a | grep -v "IMAGE ID" | awk '{print $3}')
docker rmi -f $imgs &> /dev/null

