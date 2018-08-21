#!/bin/bash

get 18edf2d18fa331265c36516a4a19ba75d26f46eafcc5e0c2d9aa6c237e8bc110 \
        "https://github.com/openresty/lua-nginx-module/archive/v$LUA_NGX_VERSION.tar.gz"

get 2a69815e4ae01aa8b170941a8e1a10b6f6a9aab699dee485d58f021dd933829a \
        "https://github.com/openresty/lua-upstream-nginx-module/archive/v$LUA_UPSTREAM_VERSION.tar.gz"

get d4a9ed0d2405f41eb0178462b398afde8599c5115dcc1ff8f60e2f34a41a4c21 \
        "https://github.com/openresty/lua-resty-lrucache/archive/v0.07.tar.gz"

get 92fd006d5ca3b3266847d33410eb280122a7f6c06334715f87acce064188a02e \
        "https://github.com/openresty/lua-resty-core/archive/v0.1.14rc1.tar.gz"

get eaf84f58b43289c1c3e0442ada9ed40406357f203adc96e2091638080cb8d361 \
        "https://github.com/openresty/lua-resty-lock/archive/v0.07.tar.gz"

get 3917d506e2d692088f7b4035c589cc32634de4ea66e40fc51259fbae43c9258d \
        "https://github.com/hamishforbes/lua-resty-iputils/archive/v0.3.0.tar.gz"

get 5d16e623d17d4f42cc64ea9cfb69ca960d313e12f5d828f785dd227cc483fcbd \
        "https://github.com/openresty/lua-resty-upload/archive/v0.10.tar.gz"

get feacc662fd7724741c2b3277b2d27b5ab2821bdb28b499d063dbd23414447249 \
        "https://github.com/openresty/lua-resty-dns/archive/v0.21rc2.tar.gz"

get 30a68f1828ed6a53ee6ed062132ea914201076058b1d126ea90ff8e55df09daf \
        "https://github.com/openresty/lua-resty-string/archive/v0.11rc1.tar.gz"

get 1ad2e34b111c802f9d0cdf019e986909123237a28c746b21295b63c9e785d9c3 \
        "http://luajit.org/download/LuaJIT-2.1.0-beta3.tar.gz"

# luajit is not available on ppc64le and s390x
if [[ (${ARCH} != "ppc64le") && (${ARCH} != "s390x") ]]; then
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

  # build and install lua-resty-waf with dependencies
  /install_lua_resty_waf.sh

fi

ln -s /usr/lib/x86_64-linux-gnu/liblua5.1.so /usr/lib/liblua.so