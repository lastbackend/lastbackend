#!/bin/bash

get_src ce66acf943a604ef9a0bb477c7efca1fe583076991647aa646aa3d8804328364 \
        "https://github.com/opentracing-contrib/nginx-opentracing/archive/v$NGINX_OPENTRACING_VERSION.tar.gz"

get_src 06dc5f9740d27dc4684399e491211be46a8069a10277f25513dadeb71199ce4c \
        "https://github.com/opentracing/opentracing-cpp/archive/v$OPENTRACING_CPP_VERSION.tar.gz"

get_src b65bb78bcd8806cf11695b980577abb5379369929240414c75eb4623a4d45cc3 \
        "https://github.com/rnburn/zipkin-cpp-opentracing/archive/v$ZIPKIN_CPP_VERSION.tar.gz"



get_src 841916d60fee16fe245b67fe6938ad861ddd3f3ecf0df561d764baeda8739362 \
        "https://github.com/jaegertracing/jaeger-client-cpp/archive/v$JAEGER_VERSION.tar.gz"


# build opentracing lib
cd "$BUILD_PATH/opentracing-cpp-$OPENTRACING_CPP_VERSION"
mkdir .build
cd .build
cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTING=OFF ..
make
make install

# build zipkin lib
cd "$BUILD_PATH/jaeger-client-cpp-$JAEGER_VERSION"
sed -i 's/-Werror//' CMakeLists.txt
mkdir .build
cd .build
cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=1 -DBUILD_TESTING=OFF -DJAEGERTRACING_WITH_YAML_CPP=OFF ..
make
make install