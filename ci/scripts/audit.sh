#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-frontend-cookie-controller
  make audit
popd 