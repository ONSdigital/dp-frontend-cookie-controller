#!/bin/bash -eux

pushd dp-frontend-cookie-controller
  make build
  cp build/dp-frontend-cookie-controller Dockerfile.concourse ../build
popd