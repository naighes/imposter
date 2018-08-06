#!/usr/bin/env bash
OWNER=naighes
PRODUCT_NAME=imposter
GIT_COMMIT=$(git rev-parse HEAD)

function get_version {
  local proj_dir=$1
  local os=$(go env GOOS)
  local arch=$(go env GOARCH)
  local raw_version=$($proj_dir/pkg/${os}_${arch}/${PWD##*/} version)
  VERSION=${raw_version#* v}
}
