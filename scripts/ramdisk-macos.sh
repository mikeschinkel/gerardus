#!/usr/bin/env bash

set -eo pipefail

# Expect "-d" for delete or "" to create
DETACH_DISK="$1"

# Number of 512 byte blocks
BLOCKS=1048576

# Mount point for ramdisk
DISK_PATH="data"

function attach_ramdisk() {
  hdiutil attach -nomount "ram://${blocks}" 2>&1 | awk '{print $1}'
}

function has_ramdisk() {
  local disk_path="$1"
  local result

  result="$(ramdisk_device "${disk_path}")"
  test "${result}" != ""
}

function ramdisk_device() {
  local disk_path="$1"

  hdiutil info | grep "${disk_path}" | awk '{print $1}'
}

function main() {
  local root_dir
  local device
  local blocks="$1"
  local disk_path="$2"
  local detach="$3"

  root_dir="$(git rev-parse --show-toplevel)"

  disk_path="${root_dir}/${disk_path}"

  if [ "${detach}" == "-d" ] ; then
    device="$(ramdisk_device "${disk_path}")"
    echo "Detaching ramdisk at ${device} from path ${disk_path}"
    hdiutil detach "${device}"
    exit $?
  fi

  if has_ramdisk "${disk_path}"; then
    echo "Ramdisk path already exists: ${disk_path}"
    exit 1
  fi

  mkdir -p "${disk_path}"

  device="$(attach_ramdisk "${blocks}")"
  # shellcheck disable=SC2181
  if [ $? -ne 0 ] ; then
    echo "ERROR: ${device}"
    exit 1
  fi

  echo "Ramdisk attached as ${device}"

  newfs_hfs -v data "${device}"
  mount -t hfs "${device}" "${disk_path}"

  test_file="${disk_path}/.test"
  touch "${test_file}"
  rm "${test_file}"

  echo "Ramdisk created at ${disk_path}"
}
main $BLOCKS "${DISK_PATH}" "${DETACH_DISK}"