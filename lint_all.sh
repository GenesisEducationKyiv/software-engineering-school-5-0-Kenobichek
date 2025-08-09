#!/usr/bin/env bash

set -e

SERVICES_DIR="internal/services"

for d in "$SERVICES_DIR"/*; do
  if [ -f "$d/go.mod" ]; then
    echo "\nLinting $d"
    (cd "$d" && golangci-lint run --fix)
  fi
done

echo -e "\nAll services linted." 