#!/usr/bin/env bash
set -eu

PROJECT_DIR=$(git rev-parse --show-toplevel)
cd $PROJECT_DIR
source ${PROJECT_DIR}/scripts/lib.sh

get_version ${PROJECT_DIR}

docker push "${OWNER}/${PRODUCT_NAME}:${VERSION}"
docker push "${OWNER}/${PRODUCT_NAME}:latest"
