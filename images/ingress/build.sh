#!/bin/sh

export NGINX_VERSION=1.13.11
export NDK_VERSION=0.3.0
export LUA_MODULE_VERSION=0.10.12rc2
export LUA_UPSTREAM_VERSION=0.07

export LUAJIT_LIB=/usr/lib
export LUAJIT_INC=/usr/include/luajit-2.0

export BUILD_PATH=/tmp/build

ARCH=$(uname -m)

if [[ ${ARCH} == "s390x" ]]; then
  git config --global pack.threads "1"
fi

addgroup -S nginx
adduser -D -S -h /var/cache/nginx -s /sbin/nologin -G nginx nginx

apk add --no-cache --virtual .build-deps \
    gcc \
		libgcc \
		libc-dev \
		make \
		openssl-dev \
		pcre-dev \
		zlib-dev \
		linux-headers \
		curl \
		gnupg \
		libxslt-dev \
		gd-dev \
		geoip-dev \
		luajit-dev

mkdir -p "$BUILD_PATH"
cd "$BUILD_PATH"

get_dep()
{
  url="$1"
  hash="$2"
  f=$(basename "$url")

  curl -sSL "$url" -o "$f"
  echo "$hash  $f" | sha256sum -c - || exit 10
  tar xzf "$f"
  rm -rf "$f"
}

# GEO IP
GEOIP_FOLDER=/etc/nginx/geoip

geoip_get ()
{
  url="$1"
  f=$(basename "$url")
  curl -sSL "$url" -o "$f"
  gunzip ${GEOIP_FOLDER}/${f}
  rm -rf "$f"
}
mkdir -p ${GEOIP_FOLDER}

cd "$GEOIP_FOLDER"
geoip_get "https://geolite.maxmind.com/download/geoip/database/GeoLiteCountry/GeoIP.dat.gz"
geoip_get "https://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz"
geoip_get "http://download.maxmind.com/download/geoip/database/asnum/GeoIPASNum.dat.gz"


cd "$BUILD_PATH"
get_dep "http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz" \
    35799c974644d2896b34ba876461dfd142c1b11f06f5aa57d255a77d4da36f05

get_dep "https://github.com/simpl/ngx_devel_kit/archive/v$NDK_VERSION.tar.gz" \
    88e05a99a8a7419066f5ae75966fb1efc409bad4522d14986da074554ae61619

get_dep "https://github.com/openresty/lua-nginx-module/archive/v$LUA_MODULE_VERSION.tar.gz" \
    18edf2d18fa331265c36516a4a19ba75d26f46eafcc5e0c2d9aa6c237e8bc110

get_dep "https://github.com/openresty/lua-upstream-nginx-module/archive/v$LUA_UPSTREAM_VERSION.tar.gz" \
    2a69815e4ae01aa8b170941a8e1a10b6f6a9aab699dee485d58f021dd933829a

get_dep "https://github.com/openresty/lua-resty-lrucache/archive/v0.07.tar.gz" \
    d4a9ed0d2405f41eb0178462b398afde8599c5115dcc1ff8f60e2f34a41a4c21

get_dep "https://github.com/openresty/lua-resty-core/archive/v0.1.14rc1.tar.gz" \
    92fd006d5ca3b3266847d33410eb280122a7f6c06334715f87acce064188a02e

get_dep "https://github.com/openresty/lua-resty-lock/archive/v0.07.tar.gz" \
    eaf84f58b43289c1c3e0442ada9ed40406357f203adc96e2091638080cb8d361

get_dep "https://github.com/hamishforbes/lua-resty-iputils/archive/v0.3.0.tar.gz" \
    3917d506e2d692088f7b4035c589cc32634de4ea66e40fc51259fbae43c9258d

get_dep "https://github.com/openresty/lua-resty-upload/archive/v0.10.tar.gz" \
    5d16e623d17d4f42cc64ea9cfb69ca960d313e12f5d828f785dd227cc483fcbd

get_dep "https://github.com/openresty/lua-resty-dns/archive/v0.21rc2.tar.gz" \
    feacc662fd7724741c2b3277b2d27b5ab2821bdb28b499d063dbd23414447249

get_dep "https://github.com/openresty/lua-resty-string/archive/v0.11rc1.tar.gz" \
    30a68f1828ed6a53ee6ed062132ea914201076058b1d126ea90ff8e55df09daf

get_dep "http://luajit.org/download/LuaJIT-2.1.0-beta3.tar.gz" \
    1ad2e34b111c802f9d0cdf019e986909123237a28c746b21295b63c9e785d9c3

# improve compilation times
CORES=$(($(grep -c ^processor /proc/cpuinfo) - 0))

export MAKEFLAGS=-j${CORES}
export CTEST_BUILD_FLAGS=${MAKEFLAGS}

# luajit is not available on ppc64le and s390x
if [[ ${ARCH} != "ppc64le" && ${ARCH} != "s390x" ]]; then

  cd "$BUILD_PATH/LuaJIT-2.1.0-beta3"
  make
  make install

  ln -sf luajit-2.1.0-beta3 /usr/local/bin/luajit

  export LUAJIT_LIB=/usr/local/lib
  export LUAJIT_INC=/usr/local/include/luajit-2.1
  export LUA_LIB_DIR="$LUAJIT_LIB/lua"

  cd "$BUILD_PATH/lua-resty-core-0.1.14rc1"
  make install

  cd "$BUILD_PATH/lua-resty-lrucache-0.07"
  make install

  cd "$BUILD_PATH/lua-resty-lock-0.07"
  make install

  cd "$BUILD_PATH/lua-resty-iputils-0.3.0"
  make install

  cd "$BUILD_PATH/lua-resty-upload-0.10"
  make install

  cd "$BUILD_PATH/lua-resty-dns-0.21rc2"
  make install

  cd "$BUILD_PATH/lua-resty-string-0.11rc1"
  make install
fi


cd "$BUILD_PATH/nginx-$NGINX_VERSION"

WITH_FLAGS="--with-debug \
  --with-http_ssl_module \
  --with-http_stub_status_module \
  --with-http_realip_module \
  --with-http_auth_request_module \
  --with-http_addition_module \
  --with-http_dav_module \
  --with-http_geoip_module \
  --with-http_gzip_static_module \
  --with-http_sub_module \
  --with-http_v2_module \
  --with-stream \
  --with-stream_ssl_module \
  --with-stream_ssl_preread_module \
  --with-threads \
  --with-http_secure_link_module"

WITH_MODULES="--add-module=$BUILD_PATH/ngx_devel_kit-$NDK_VERSION \
  --add-module=$BUILD_PATH/lua-nginx-module-$LUA_MODULE_VERSION \
  --add-module=$BUILD_PATH/lua-upstream-nginx-module-$LUA_UPSTREAM_VERSION"

./configure \
  --prefix=/etc/nginx \
  --sbin-path=/usr/sbin/nginx \
  --conf-path=/etc/nginx/nginx.conf \
  --modules-path=/etc/nginx/modules \
  --http-log-path=/var/log/nginx/access.log \
  --error-log-path=/var/log/nginx/error.log \
  --lock-path=/var/lock/nginx.lock \
  --pid-path=/run/nginx.pid \
  --http-client-body-temp-path=/var/lib/nginx/body \
  --http-fastcgi-temp-path=/var/lib/nginx/fastcgi \
  --http-proxy-temp-path=/var/lib/nginx/proxy \
  --http-scgi-temp-path=/var/lib/nginx/scgi \
  --http-uwsgi-temp-path=/var/lib/nginx/uwsgi \
  ${WITH_FLAGS} \
  --without-mail_pop3_module \
  --without-mail_smtp_module \
  --without-mail_imap_module \
  --without-http_uwsgi_module \
  --without-http_scgi_module \
  ${WITH_MODULES} \
  && make || exit 1 \
  && make install || exit 1

mv /usr/share/nginx/sbin/nginx /usr/sbin

apk del .build-deps

