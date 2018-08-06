#!/usr/bin/env bash
RELEASE=0
while [[ $# -gt 0 ]] && [[ ."$1" = .--* ]] ;
do
  case $1 in
    --release)
      shift
      RELEASE=1
      ;;
    *)
      echo "parameter $1 is not recognized as a valid option"
      exit 1
  esac
done

PROJECT_DIR=$(git rev-parse --show-toplevel)
cd $PROJECT_DIR
source ${PROJECT_DIR}/scripts/lib.sh

XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-linux darwin windows freebsd openbsd solaris}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"

rm -rf ${PROJECT_DIR}/pkg/*

if ! which gox > /dev/null; then
  go get -u github.com/mitchellh/gox
fi

export CGO_ENABLED=0

LD_FLAGS="-X main.GitCommit=${GIT_COMMIT} $LD_FLAGS"

if (( $RELEASE == 1 )) ; then
  LD_FLAGS="-X main.GitCommit=${GIT_COMMIT} -X github.com/${OWNER}/${PRODUCT_NAME}/version.Prerelease= -s -w"
fi

go get -d

gox \
  -os="${XC_OS}" \
  -arch="${XC_ARCH}" \
  -osarch="${XC_EXCLUDE_OSARCH}" \
  -ldflags "${LD_FLAGS}" \
  -output "${PROJECT_DIR}/pkg/{{.OS}}_{{.Arch}}/${PRODUCT_NAME}"

for PLATFORM in $(find ${PROJECT_DIR}/pkg -mindepth 1 -maxdepth 1 -type d); do
  OSARCH=$(basename ${PLATFORM})
  pushd $PLATFORM >/dev/null 2>&1
  zip ../${OSARCH}.zip ./*
  popd >/dev/null 2>&1
done
