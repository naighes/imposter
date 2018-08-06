#!/usr/bin/env bash
set -e

GH_AUTH_TOKEN=$1
shift

if [[ -z $GH_AUTH_TOKEN ]]; then
  echo "a token is required for releasing"
  exit 1
fi

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

GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/$OWNER/$PRODUCT_NAME"
GH_AUTH_HEADER="Authorization: token $GH_AUTH_TOKEN"
echo "releasing '$RAW_VERSION'"
curl -s \
  -o /dev/null \
  -H "$GH_AUTH_HEADER" \
  $GH_REPO || { echo "invalid repo, token or network issue";  exit 1; }
PAYLOAD="{\"tag_name\":\"$VERSION\",\"target_commitish\":\"master\",\"name\":\"$RAW_VERSION\",\"body\":\"\",\"draft\":false,\"prerelease\":false}"
# TODO: check if a release with that tag has been already released
BASE_URL=$(curl -s \
  -X POST \
  -H "$GH_AUTH_HEADER" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD" \
  ${GH_REPO}/releases | jq -r .upload_url) || { echo "invalid repo, token or network issue";  exit 1; }
for FILENAME in $(find ./pkg/dist -mindepth 1 -maxdepth 1 -type f); do
  GH_ASSET="${BASE_URL%\{*}?name=$(basename $FILENAME)"
  echo "publishing asset '$GH_ASSET'"
  curl -s \
    --data-binary @"$FILENAME" \
    -o /dev/null \
    -H "$GH_AUTH_HEADER" \
    -H "Content-Type: application/octet-stream" \
    $GH_ASSET || { echo "invalid repo, token or network issue";  exit 1; }
done
${PROJECT_DIR}/scripts/docker/push.sh

echo "done"
exit 0
