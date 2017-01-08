FROM alpine:3.4

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin:/usr/local/bin

RUN apk add --update --no-cache \
      ca-certificates

RUN apk add --update --no-cache --virtual build-dependencies \
      git \
      go \
      gcc \
      alpine-sdk

ADD . /go/src/github.com/lastbackend/lastbackend
WORKDIR /go/src/github.com/lastbackend/lastbackend

RUN go get github.com/tools/godep \
    && godep restore \
    && make build \
    && make install \
    && apk del --purge build-dependencies

WORKDIR /go/bin
RUN rm -rf /go/pkg \
    && rm -rf /go/src \
    && rm -rf /var/cache/apk/*

EXPOSE 2967

CMD ["/usr/local/bin/lastbackend", "daemon", "-c", "/etc/lastbackend/config.yml"]
