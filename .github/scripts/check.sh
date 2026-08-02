#!/usr/bin/env bash

set -euo pipefail

action_directories=(
  ".github/actions/get-dockerhub-version-tag"
  ".github/actions/get-pecl-package-version"
)

for directory in "${action_directories[@]}"; do
  echo "Testing ${directory}"
  (
    cd "${directory}"
    go test ./...
    go vet ./...
  )

  echo "Building ${directory}"
  docker build "${directory}"
done

echo "Building WordPress image"
docker build .
