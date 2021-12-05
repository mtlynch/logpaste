#!/bin/bash

set -e

SFTP_HOST="$1"
shift
SFTP_PORT="$1"
shift

timeout 15 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/$0/$1; do sleep 1; done' "${SFTP_HOST}" "${SFTP_PORT}"

>&2 echo "sftp is up - executing command"

exec "$@"
