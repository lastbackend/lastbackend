#!/usr/bin/env bash

WORKDIR=${TRAVIS_BUILD_DIR:-`pwd`}
TAG=latest
SUBDIR=""

echo "Deleting old output"
rm -rf ${WORKDIR}/docs/output
mkdir -p ${WORKDIR}/docs/output/${SUBDIR}

echo "Copying images"
cp -R ${WORKDIR}/docs/assets ${WORKDIR}/docs/output/${SUBDIR}assets

echo "Copying files"
cp -R ${WORKDIR}/docs/files/* ${WORKDIR}/docs/output/

echo "Generating docs"
docker run -v ${WORKDIR}/docs/:/documents/ --name asciidoc-to-html --rm asciidoctor/docker-asciidoctor asciidoctor -a revnumber=${TAG} -D /documents/output/${SUBDIR} index.adoc

