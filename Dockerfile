FROM golang:1.11

LABEL maintainer="Tony Lee <tungsheng@gmail.com>"

ADD . /go/src

EXPOSE 3003 9091

WORKDIR /go/src

RUN go get -u github.com/ipfs/go-ipfs
RUN go get -u github.com/joho/godotenv
RUN go get -u github.com/rs/zerolog/log
RUN go get -u github.com/segmentio/go-env
RUN go get -u gopkg.in/go-playground/webhooks.v5/bitbucket
RUN go run main.go
