FROM golang:1.9.2-alpine3.6

WORKDIR /go/src/app

COPY . .

RUN apk add --no-cache git

RUN go-wrapper download

RUN go-wrapper install

CMD ["go-wrapper", "run"]

EXPOSE 8080
