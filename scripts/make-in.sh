#!/usr/bin/env bash

wd_file="$1/MAKE_DIR"
if [ ! -f "${wd_file}" ]; then
  echo "File ${wd_file} does not exist."
  echo "   It should contain current working directory for Makefile."
fi
cd "$(cat "${wd_file}")" || echo "Cannot cd into ${wd_file}"
make "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9"






