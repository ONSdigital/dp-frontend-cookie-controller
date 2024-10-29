#!/bin/bash -eux

pushd dp-frontend-cookie-controller
DEFAULT_WORKSPACE=$(pwd)
export DEFAULT_WORKSPACE
bash -ex /entrypoint.sh
popd
