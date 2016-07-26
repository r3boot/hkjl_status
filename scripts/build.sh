#!/bin/bash

SETTINGS="$(dirname ${0})/settings.sh"

source ${SETTINGS}

echo ">>> Setting up build environment"
$(dirname ${0})/cleanup.sh
mkdir -v ${BUILD_DIR}

GOPATH="$(pwd)"
export GOPATH

echo ">>> Fetching dependencies"
${GO} get -v github.com/r3boot/rlib

echo ">>> Building binary"
${GO} build -v -o ${BUILD_DIR}/hkjl_status

echo ">>> Stripping debug symbols"
${STRIP} ${BUILD_DIR}/hkjl_status

stat ${BUILD_DIR}/hkjl_status
