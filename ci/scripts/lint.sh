#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-frontend-cookie-controller
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
  make lint
popd
