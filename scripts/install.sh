#!/bin/bash

################################################################################
#                                                                              #
#                                es Installer                                  #
#                                                                              #
################################################################################

# This script will install the latest version of es on your machine from
# the precompiled binary releases in the official repository.

# check platform
if [[ "$(uname)" == "Darwin" ]]; then
  PLATFORM="darwin"
elif [[ "$(uname)" == "Linux" ]]; then
  PLATFORM="linux"
else
  PLATFORM="windows"
fi

install_bin() {
  local name=$1
  if [[ -z "$name" ]]; then
    echo "install_bin: name is empty"
    return 1
  fi
  echo "install_bin: installing $name"
  LATEST_DOWNLOAD_PREFIX="https://github.com/robertlestak/es/releases/latest/download/"
  FILE_NAME="${name}_${PLATFORM}"
  DL="${LATEST_DOWNLOAD_PREFIX}${FILE_NAME}"
  echo "install_bin: downloading $DL"
  CURL -s -L $DL > $FILE_NAME
  chmod +x $FILE_NAME
  mv $FILE_NAME /usr/local/bin/$name
}

main() {
  install_bin "es"
}
main "$@"