#!/usr/bin/env bash
set -eu

PROJECT_DIR=$(git rev-parse --show-toplevel)
cd $PROJECT_DIR
source ${PROJECT_DIR}/scripts/lib.sh

get_version ${PROJECT_DIR}

docker build \
    --no-cache \
    --build-arg=PRODUCT_NAME=${PRODUCT_NAME} \
    --build-arg=VERSION=${VERSION} \
    --build-arg=OWNER=${OWNER} \
    -t ${OWNER}/$PRODUCT_NAME:$VERSION \
    .

docker tag \
    "${OWNER}/${PRODUCT_NAME}:${VERSION}" \
    "${OWNER}/${PRODUCT_NAME}:latest"
