FROM golang:1-alpine AS builder

RUN mkdir -p /go/src/github.com/ihciah/rabbit-tcp
COPY . /go/src/github.com/ihciah/rabbit-tcp

RUN apk upgrade \
    && apk add git \
    && apk add make \
    && cd /go/src/github.com/ihciah/rabbit-tcp \
    && go get -v -t -d ./... \
    && make

FROM alpine:latest AS dist
LABEL maintainer="ihciah <ihciah@gmail.com>"

ENV MODE s
ENV PASSWORD PASSWORD
ENV RABBITADDR :443
ENV LISTEN :9891
ENV DEST=
ENV TUNNELN 6
ENV VERBOSE 2

COPY --from=builder /go/src/github.com/ihciah/rabbit-tcp/bin/rabbit /usr/bin/rabbit

CMD exec rabbit \
      --mode=$MODE \
      --password=$PASSWORD \
      --rabbit-addr=$RABBITADDR \
      --listen=$LISTEN \
      --dest=$DEST \
      --tunnelN=$TUNNELN \
      --verbose=$VERBOSE