#!/usr/bin/env bash

set -e

go mod tidy

bazel run //:gazelle-update-repos

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Please run 'bazel run //:gazelle-update-repos' and commit the changes"
  exit 1
fi

go fmt ./...

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Please run 'go fmt ./...' and commit the changes"
  exit 1
fi