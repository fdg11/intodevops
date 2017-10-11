#!/usr/bin/env bash

if [ $# != "2" ]; then

echo -e "Enter bot token and chat id"  
exit 1 

fi 

cat <<EOF > $(pwd)/golang/Dockerfile
FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go-wrapper download

RUN go-wrapper install

CMD ["go-wrapper", "run", "-telegrambottoken", "$1", "-chatid", "$2"]

EXPOSE 8080
EOF

echo -e "Dockerfile create with credentials"
