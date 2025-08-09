#!/usr/bin/env bash

set -e

SERVICES_DIR="internal/services"

for d in "$SERVICES_DIR"/*; do
  if [ -f "$d/go.mod" ]; then
    echo -e "\nRunning tests for $d"
    (cd "$d" && go test -v ./...)
  fi
done

echo -e "\nAll tests completed." 