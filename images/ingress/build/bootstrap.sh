#!/bin/bash

bootstrap () {
apt-get update && apt-get dist-upgrade -y

# install required packages to build
clean-install \
  bash \
  build-essential \
  curl ca-certificates \
  libgeoip1 \
  libgeoip-dev \
  patch \
  libpcre3 \
  libpcre3-dev \
  libssl-dev \
  zlib1g \
  zlib1g-dev \
  libaio1 \
  libaio-dev \
  openssl \
  libperl-dev \
  cmake \
  util-linux \
  lua5.1 liblua5.1-0 liblua5.1-dev \
  lmdb-utils \
  libjemalloc1 libjemalloc-dev \
  wget \
  libcurl4-openssl-dev \
  procps \
  git g++ pkgconf flex bison doxygen libyajl-dev liblmdb-dev libtool dh-autoreconf libxml2 libpcre++-dev libxml2-dev \
  lua-cjson \
  python \
  luarocks \
  || exit 1
}

clean () {

  apt-mark unmarkauto \
    bash \
    curl ca-certificates \
    libgeoip1 \
    libpcre3 \
    zlib1g \
    libaio1 \
    xz-utils \
    geoip-bin \
    libyajl2 liblmdb0 libxml2 libpcre++ \
    gzip \
    openssl

  apt-get remove -y --purge \
    build-essential \
    gcc-6 \
    cpp-6 \
    libgeoip-dev \
    libpcre3-dev \
    libssl-dev \
    zlib1g-dev \
    libaio-dev \
    linux-libc-dev \
    cmake \
    wget \
    git g++ pkgconf flex bison doxygen libyajl-dev liblmdb-dev libgeoip-dev libtool dh-autoreconf libpcre++-dev libxml2-dev

  apt-get autoremove -y

  mkdir -p /var/lib/nginx/body /usr/share/nginx/html

  mv /usr/share/nginx/sbin/nginx /usr/sbin

  cd /

  rm -rf "$BUILD_PATH"
  rm -Rf /usr/share/man /usr/share/doc
  rm -rf /tmp/* /var/tmp/*
  rm -rf /var/lib/apt/lists/*
  rm -rf /var/cache/apt/archives/*
  rm -rf /usr/local/modsecurity/bin
  rm -rf /usr/local/modsecurity/include
  rm -rf /usr/local/modsecurity/lib/libmodsecurity.a

  rm -rf /etc/nginx/owasp-modsecurity-crs/.git
  rm -rf /etc/nginx/owasp-modsecurity-crs/util/regression-tests
}