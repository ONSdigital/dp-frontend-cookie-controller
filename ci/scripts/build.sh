#!/bin/bash -eux

cwd=$(pwd)

pushd dp-frontend-cookie-controller
  make build && cp build/dp-frontend-cookie-controller $cwd/build
  cp Dockerfile.concourse $cwd/build
popd