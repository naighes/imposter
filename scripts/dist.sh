#!/usr/bin/env bash
set -e

RELEASE=0
while [[ $# -gt 0 ]] && [[ ."$1" = .--* ]] ;
do
  case $1 in
    --release)
      shift
      RELEASE=1
      AUTH_TOKEN=$1
      if [[ -z $AUTH_TOKEN ]]; then
        echo "a token is required in release mode"
        exit 1
      fi
      shift
      ;;
    *)
      echo "parameter $1 is not recognized as a valid option"
      exit 1
  esac
done

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd $DIR

XC_OS=$(go env GOOS)
XC_ARCH=$(go env GOARCH)
RAW_VERSION=$(./pkg/${XC_OS}_${XC_ARCH}/${PWD##*/} version)
VERSION=${RAW_VERSION#* v}
PRODUCT_NAME=$(echo "${RAW_VERSION% v*}" | tr '[:upper:]' '[:lower:]')

echo "cleaning"
rm -rf ./pkg/dist
mkdir -p ./pkg/dist

for FILENAME in $(find ./pkg -mindepth 1 -maxdepth 1 -type f); do
  FILENAME=$(basename $FILENAME)
  SOURCE_FILE=./pkg/${FILENAME}
  TARGET_FILE=./pkg/dist/${PRODUCT_NAME}_${VERSION}_${FILENAME}
  echo "copying '$SOURCE_FILE' to '$TARGET_FILE'"
  cp $SOURCE_FILE $TARGET_FILE
done

if (( $RELEASE == 1 )) ; then
  OWNER="naighes"
  GH_API="https://api.github.com"
  GH_REPO="$GH_API/repos/$OWNER/$PRODUCT_NAME"
  AUTH_HEADER="Authorization: token $AUTH_TOKEN"
  echo "releasing '$RAW_VERSION'"
  curl -s \
    -o /dev/null \
    -H "$AUTH_HEADER" \
    $GH_REPO || { echo "invalid repo, token or network issue";  exit 1; }
  PAYLOAD="{\"tag_name\":\"$VERSION\",\"target_commitish\":\"master\",\"name\":\"$RAW_VERSION\",\"body\":\"\",\"draft\":false,\"prerelease\":false}"
  # TODO: check if a release with that tag has been already released
  BASE_URL=$(curl -s \
    -X POST \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD" \
    ${GH_REPO}/releases | jq -r .upload_url) || { echo "invalid repo, token or network issue";  exit 1; }
  for FILENAME in $(find ./pkg/dist -mindepth 1 -maxdepth 1 -type f); do
    GH_ASSET="${BASE_URL%\{*}?name=$(basename $FILENAME)"
    echo "publishing asset '$GH_ASSET'"
    curl -s \
      --data-binary @"$FILENAME" \
      -o /dev/null \
      -H "$AUTH_HEADER" \
      -H "Content-Type: application/octet-stream" \
      $GH_ASSET || { echo "invalid repo, token or network issue";  exit 1; }
  done
fi
echo "done"
exit 0
