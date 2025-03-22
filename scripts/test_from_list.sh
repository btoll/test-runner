#!/bin/bash

FILE=${1:-operators.list}

#        --fail-fast \
#        --fail-on-empty \

while read -r op
do
    OPERATOR=$(awk -F '/' '{print $2}' <<< "$op")
    DISABLE_JUNIT_REPORT=true ginkgo run \
        --tags=e2e,osde2e \
        --flake-attempts=3 \
        --trace \
        --output-dir tests \
        --json-report "$OPERATOR.json" \
        "$op"
done < "$FILE"

