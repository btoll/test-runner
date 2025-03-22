#!/bin/bash
# shellcheck disable=2045

jprint() {
    local report="$1"
    local test="$2"
    local state="$3"

    blob=$(printf "%s" "$report" | jq -r '.SpecReports[] | select(.State=="'"$state"'")')
    if [ -n "$blob" ]
    then
        printf "Test file: %s\n" "tests/$test"
        printf "%s" "$report" | jq -r '.PreRunStats | to_entries[] | "\(.key): \(.value)"'
        leaf_node_text=$(printf "%s" "$blob" | jq -r '.LeafNodeText')
        if [ -n "$leaf_node_text" ]
        then
            # $state will be uppercased.
            printf "%s test name: %s\n" "${state^^}" "$leaf_node_text"
        fi
        printf "\n"
    fi
}

for test in $(ls tests)
do
    report=$(< "tests/$test" jq -r '.[]')
    jprint "$report" "$test" failed
    jprint "$report" "$test" pending
    jprint "$report" "$test" skipped
done

