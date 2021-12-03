#!/bin/bash

pid=0

# SIGTERM-handler
term_handler() {
  if [ $pid -ne 0 ]; then
    kill -SIGTERM "$pid"
    wait "$pid"
  fi

  exit
}

trap 'kill ${!}; term_handler' SIGTERM SIGINT

if [ ! -z "${PRELOAD_EXTENSIONS}" ]; then
    ext-preloader -ext "${PRELOAD_EXTENSIONS}" &
    pid="$!"
fi

"$@"

if [ $pid -ne 0 ]; then
  wait "$pid"
fi
