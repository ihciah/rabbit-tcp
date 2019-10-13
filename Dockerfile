FROM golang:1-alpine AS builder

RUN mkdir -p /go/src/github.com/ihciah/rabbit-tcp
COPY . /go/src/github.com/ihciah/rabbit-tcp

RUN apk upgrade \
    && apk add git \
    && go get github.com/ihciah/rabbit-tcp/cmd

FROM alpine:latest AS dist
LABEL maintainer="ihciah <ihciah@gmail.com>"

ENV MODE s
ENV PASSWORD PASSWORD
ENV RABBITADDR :443
ENV LISTEN :9891
ENV DEST=
ENV TUNNELN 6
ENV VERBOSE 2

COPY --from=builder /go/bin/cmd /usr/bin/rabbit

CMD exec rabbit \
      --mode=$MODE \
      --password=$PASSWORD \
      --rabbit-addr=$RABBITADDR \
      --listen=$LISTEN \
      --dest=$DEST \
      --tunnelN=$TUNNELN \
      --verbose=$VERBOSE