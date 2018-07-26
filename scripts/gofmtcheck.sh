#!/usr/bin/env bash
gofmt_files=$(gofmt -l `find . -name '*.go'`)
if [[ -n ${gofmt_files} ]]; then
  echo 'gofmt needs running on the following files:'
  echo "${gofmt_files}"
  exit 1
fi
exit 0
