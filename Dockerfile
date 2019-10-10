FROM golang:1.12.7-alpine3.10 AS builder

RUN mkdir -p /go/src/github.com/ihciah/rabbit-tcp
COPY . /go/src/github.com/ihciah/rabbit-tcp

RUN apk upgrade \
    && apk add git \
    && go get github.com/ihciah/rabbit-tcp/cmd

FROM alpine:3.10 AS dist
LABEL maintainer="ihciah <ihciah@gmail.com>"

ENV MODE s
ENV PASSWORD PASSWORD
ENV RABBITADDR :443
ENV LISTEN :9891
ENV DEST=
ENV TUNNELN 6
ENV VERBOSE 2

RUN apk upgrade \
    && apk add tzdata \
    && rm -rf /var/cache/apk/*

COPY --from=builder /go/bin/cmd /usr/bin/rabbit

CMD exec rabbit \
      --mode=$MODE \
      --password=$PASSWORD \
      --rabbit-addr=$RABBITADDR \
      --listen=$LISTEN \
      --dest=$DEST \
      --tunnelN=$TUNNELN \
      --verbose=$VERBOSE