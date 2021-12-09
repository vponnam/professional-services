#!/bin/bash

set -eo pipefail

name=$1
project=$2
secret=$3

rotate_secret() {
    for value in $(gcloud secrets list  --format json  | jq -r '.[].name' | cut -d "/" -f 4); do
        if [[ "${value}" == "${name}" ]]; then
            # Delete if the key exists
            destroy_secret
            
            break
        fi
    done
    create_secret
}

create_secret() {
    printf "${secret}" | gcloud secrets create "${name}" --project "${project}" --data-file=-
}

destroy_secret() {
    gcloud secrets delete "${name}" --project "${project}" -q
}

rotate_secret