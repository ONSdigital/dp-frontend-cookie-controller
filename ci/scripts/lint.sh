#!/bin/bash -eux

pushd dp-frontend-cookie-controller
make prepare-lint-go
DEFAULT_WORKSPACE=$(pwd)
export DEFAULT_WORKSPACE
bash -ex /entrypoint.sh
make post-lint-go
popd
