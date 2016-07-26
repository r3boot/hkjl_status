#!/bin/bash

SETTINGS="$(dirname ${0})/settings.sh"

source ${SETTINGS}

echo ">>> Setting up build environment"
$(dirname ${0})/cleanup.sh
${GIT} checkout -b ${BRANCH}

GOPATH="$(pwd)"
export GOPATH

echo ">>> Fetching dependencies"
${GO} get -v github.com/r3boot/rlib

echo ">>> Building binary"
${GO} build -v

echo ">>> Stripping debug symbols"
${STRIP} hkjl_status

stat hkjl_status
