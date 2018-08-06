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

PROJECT_DIR=$(git rev-parse --show-toplevel)
cd $PROJECT_DIR
source ${PROJECT_DIR}/scripts/lib.sh

get_version ${PROJECT_DIR}

rm -rf ${PROJECT_DIR}/pkg/dist
mkdir -p ${PROJECT_DIR}/pkg/dist

for FILENAME in $(find ./pkg -mindepth 1 -maxdepth 1 -type f); do
  FILENAME=$(basename $FILENAME)
  SOURCE_FILE=${PROJECT_DIR}/pkg/${FILENAME}
  TARGET_FILE=${PROJECT_DIR}/pkg/dist/${PRODUCT_NAME}_${VERSION}_${FILENAME}
  echo "copying '$SOURCE_FILE' to '$TARGET_FILE'"
  cp $SOURCE_FILE $TARGET_FILE
done

${PROJECT_DIR}/scripts/docker/build.sh

if (( $RELEASE == 1 )) ; then
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
  ${PROJECT_DIR}/scripts/docker/push.sh
fi
echo "done"
exit 0
