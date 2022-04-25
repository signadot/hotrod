#!/usr/bin/env bash

set -e

if [ $# -ne 1 ]; then
        echo "usage: <namespace>"
        exit 1
fi
NAMESPACE=$1

kubectl -n ${NAMESPACE} delete -f k8s/pieces/

kubectl -n sd-mysql delete -f k8s/pieces/db-secret.yaml
