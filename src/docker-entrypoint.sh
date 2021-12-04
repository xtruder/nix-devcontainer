#!/bin/bash

if [ ! -z "${PRELOAD_EXTENSIONS}" ]; then
  ext-preloader &
fi

if [ $# -eq 0 ]; then
  while :; do sleep 2073600; done
else
  "$@" &
fi

wait -n
