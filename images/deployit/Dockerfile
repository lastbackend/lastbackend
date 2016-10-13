FROM alpine:3.4

ENV GOROOT=/usr/lib/go \
    GOPATH=/go \
    GOBIN=/go/bin \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin:/usr/local/bin

RUN apk add --update --no-cache \
      ca-certificates

RUN apk add --update --no-cache --virtual build-dependencies \
      git \
      go \
      gcc \
      alpine-sdk

WORKDIR /go/src/github.com/deployithq/deployit
ADD . /go/src/github.com/deployithq/deployit

RUN go get \
    && make build \
    && make install \
    && apk del --purge build-dependencies

WORKDIR /go/bin
RUN rm -rf /go/pkg \
    && rm -rf /go/src \
    && rm -rf /var/cache/apk/*

EXPOSE 2967

CMD ["/usr/local/bin/deployit", "daemon", "-c", "/etc/deployit/config.yml"]
