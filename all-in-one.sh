#!/bin/sh
rm -f k8s/all-in-one/demo.yaml
for f in k8s/pieces/*.yaml; do
	echo $f
	cat $f >> k8s/all-in-one/demo.yaml
	echo "---" >> k8s/all-in-one/demo.yaml
done
