#!/usr/bin/env bash

_with_arguments () {
    count=${1};
    shift;
    if [[ "$#" -lt ${count} ]]; then
        echo "missing arguments, expected at least ${count} but received $#";
        exit 1;
    fi
}

function template() {
  _with_arguments 2 "$@"
  username=${1}
  project=${2}

  files=("go.mod" "README.md" ".golangci.yaml")

  for file in "${files[@]}"; do
    perl -pi -e "s/viqueen/${username}/g" "${file}"
    perl -pi -e "s/go-template/${project}/g" "${file}"
  done
}

# shellcheck disable=SC2294
eval "$@"