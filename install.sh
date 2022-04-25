#!/usr/bin/env bash

set -e

if [ $# -ne 1 ]; then
        echo "usage: <namespace>"
        exit 1
fi
NAMESPACE=$1

kubectl create ns ${NAMESPACE}

#
# grant access to sd-mysql-8.0-plugin
#
kubectl -n sd-mysql apply -f k8s/pieces/db-secret.yaml

#
# set up pieces
#
kubectl -n ${NAMESPACE} apply -f k8s/pieces/

#
# install ref db
#
kubectl -n ${NAMESPACE} wait pod/customer-db-0 --for=condition=Ready --timeout=60s
sleep 2

kubectl -n ${NAMESPACE} port-forward customer-db-0 3307:3306 &
pfpid=$!
set +e
for i in $(seq 1 100)
do
	mysql -u root --password=abc --port 3307 < schema/customers.sql
	if [ $? -eq 0 ]; then
		echo created schema on customer-db
		break
	fi
	if [ $i -eq 100 ]; then
		echo "givng up."
		kill $pfpid
		exit 1
	fi
	echo retrying db setup
	sleep 1
done
kill $pfpid
