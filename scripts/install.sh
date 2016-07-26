#!/bin/bash

SETTINGS="$(dirname ${0})/settings.sh"

source ${SETTINGS}

if [[ ! -f "${BUILD_DIR}/hkjl_status" ]]; then
    $(dirname ${0})/build.sh
fi

getent group ${GROUP} >/dev/null
if [[ ${?} -ne 0 ]]; then
    echo ">>> Adding group '${GROUP}'"
    groupadd -r ${GROUP}
fi

getent passwd ${USER} >/dev/null
if [[ ${?} -ne 0 ]]; then
    echo ">>> Adding user '${USER}'"
    useradd -r -d "${OUTPUT_DIR}" -g ${GROUP} -G ${WWW_GROUP} -s /usr/sbin/nologin ${USER}
fi

install -v -o root -g root -m 0755 ${BUILD_DIR}/hkjl_status /usr/bin/hkjl_status
install -d -v -o ${USER} -g ${GROUP} -m 0755 ${OUTPUT_DIR}
install -d -v -o ${USER} -g ${USER} -m 0755 assets/css ${OUTPUT_DIR}/css
install -v -o ${USER} -g ${USER} -m 0644 assets/css/* ${OUTPUT_DIR}/css
install -d -v -o ${USER} -g ${USER} -m 0755 assets/imgs ${OUTPUT_DIR}/imgs
install -v -o ${USER} -g ${USER} -m 0644 assets/imgs/* ${OUTPUT_DIR}/imgs
install -d -v -o ${USER} -g ${USER} -m 0755 assets/js ${OUTPUT_DIR}/js
install -v -o ${USER} -g ${USER} -m 0644 assets/js/* ${OUTPUT_DIR}/js
install -d -v -o root -g root ${TEMPLATE_DIR}
install -v  -o root -g root templates/index.html ${TEMPLATE_DIR}/index.html

if [[ ! -f "${CRONJOB}" ]]; then
    echo ">>> Adding cronjob"
    echo -e "# hkjl_status-update-cronjob\n* * * * * ${USER} /usr/bin/hkjl_status -o \"${OUTPUT_DIR}\" -t \"${TEMPLATE_DIR}\"" > ${CRONJOB}
fi
