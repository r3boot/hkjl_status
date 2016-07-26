# File containing various settings used throughout the installation scripts

# Where to install the binary
BINDIR='/usr/bin'

# Path to output directory
OUTPUT_DIR='/srv/www/hkjl_status/htdocs'

# Path to template directory
TEMPLATE_DIR='/usr/share/hkjl_status/templates'

# File containing cronjob
CRONJOB='/etc/cron.d/update_hkjl_status'

# Path to build directory
BUILD_DIR='./build'

# User/group under which the binary should run
USER='_hkjl'
GROUP='_hkjl'
WWW_GROUP='www-data'

# Binaries used in these scripts
GIT="$(which git)"
GO="$(which go)"
STRIP="$(which strip)"

# Can only run these scripts as root
if [[ $(whoami) != "root" ]]; then
    echo "Can only run this script as root"
    exit 1
fi

# Sanity checks
if [[ -z "${GIT}" ]]; then
    echo "Git binary not found, please install the 'git' package"
    exit 1
fi

if [[ -z "${GO}" ]]; then
    echo "Go binary not found, please install the 'golang' package"
    exit 1
fi

if [[ -z "${STRIP}" ]]; then
    echo "Strip binary not found, please install the 'binutils' package"
    exit 1
fi

# Remove the lines below once you have configured this file
echo "Please configure scripts/settings.sh before running this script"
exit 1
