#!/bin/bash

SETTINGS="$(dirname ${0})/settings.sh"

source ${SETTINGS}

${GIT} commit -a -m 'Result of build'
${GIT} checkout master
${GIT} branch | grep -q ${BRANCH}
if [[ ${?} -eq 0 ]]; then
    ${GIT} branch -D ${BRANCH}
fi
