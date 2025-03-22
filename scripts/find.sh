#!/bin/bash

# This little script is executed as part of the unix `find` command:
# For example:
#   find repos/ -maxdepth 3 -type d -ipath "*/test/e2e" -exec bash scripts/find.sh {} \;

#    --fail-fast \
#    --fail-on-empty \

OPERATOR=$(awk -F '/' '{print $2}' <<< "$1")
DISABLE_JUNIT_REPORT=true ginkgo run \
    --tags=e2e,osde2e \
    --flake-attempts=3 \
    --trace \
    --output-dir tests \
    --json-report "$OPERATOR.json" \
    "$1"

