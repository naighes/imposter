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

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"

PRODUCT_NAME=${PWD##*/}
GIT_COMMIT=$(git rev-parse HEAD)
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-linux darwin windows freebsd openbsd solaris}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"

rm -rf pkg/*

if ! which gox > /dev/null; then
  go get -u github.com/mitchellh/gox
fi

export CGO_ENABLED=0

LD_FLAGS="-X main.GitCommit=${GIT_COMMIT} $LD_FLAGS"

if (( $RELEASE == 1 )) ; then
  LD_FLAGS="-X main.GitCommit=${GIT_COMMIT} -X github.com/naighes/imposter/version.Prerelease= -s -w"
fi

gox \
  -os="${XC_OS}" \
  -arch="${XC_ARCH}" \
  -osarch="${XC_EXCLUDE_OSARCH}" \
  -ldflags "${LD_FLAGS}" \
  -output "pkg/{{.OS}}_{{.Arch}}/${PRODUCT_NAME}"

for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
  OSARCH=$(basename ${PLATFORM})
  pushd $PLATFORM >/dev/null 2>&1
  zip ../${OSARCH}.zip ./*
  popd >/dev/null 2>&1
done
