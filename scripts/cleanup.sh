#!/bin/bash

SETTINGS="$(dirname ${0})/settings.sh"

source ${SETTINGS}

if [[ -d "${BUILD_DIR}" ]]; then
    rm -vrf "${BUILD_DIR}"
fi
