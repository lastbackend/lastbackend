#!/usr/bin/env bash

WORKDIR=${TRAVIS_BUILD_DIR:-`pwd`}
echo "Deleting old output"
rm -rf ${WORKDIR}/docs/output
mkdir -p ${WORKDIR}/docs/output/latest

echo "Copying images"
cp -R ${WORKDIR}/docs/assets ${WORKDIR}/docs/output/latest/assets
echo "Generating docs"
docker run -v ${WORKDIR}/docs/:/documents/ --name asciidoc-to-html --rm asciidoctor/docker-asciidoctor asciidoctor -a revnumber=latest -D /documents/output/latest index.adoc

