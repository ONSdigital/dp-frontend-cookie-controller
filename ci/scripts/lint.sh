#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-frontend-cookie-controller
  make lint
popd
