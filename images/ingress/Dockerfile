FROM golang:1.10.1-alpine3.7

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache iptables \
    linux-headers \
    gcc \
    musl-dev \
    iproute2 \
    haproxy

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
    bash \
    git  \
    make \
	\
	&& rm -rf /*.patch

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY ./images/ingress/errors /var/run/html/errors

ADD . /go/src/github.com/lastbackend/lastbackend
WORKDIR /go/src/github.com/lastbackend/lastbackend

RUN make APP=ingress build && make APP=ingress install
RUN apk del --purge .build-deps

WORKDIR /go/bin
RUN rm -rf /go/pkg \
    && rm -rf /go/src \
    && rm -rf /var/cache/apk/*

CMD ["/usr/local/bin/ingress"]