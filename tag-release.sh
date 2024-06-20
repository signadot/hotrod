#!/bin/sh
#

set -ex


RELEASE_TAG=$1

if [ -z ${RELEASE_TAG} ] ; then
	echo "usage $0 <release-tag>"
	exit 1
fi

if [ ! -d k8s/base ]; then
	echo wrong dir
	exit 1
fi

#
# check for changes before editing images
#
set +e
git diff-index  --exit-code HEAD
if [ ! $? -eq 0 ]; then
	echo "repository is dirty"
	exit 1
fi
set -e

(cd k8s/base && kustomize edit set image signadot/hotrod:${RELEASE_TAG})

# do we need to commit the kustomize changes from above?
set +e
git diff-index  --exit-code HEAD
diffCode=$?;

set -e
if [ ${diffCode} -eq 1 ] ; then 
	git add k8s/base
	git status
	git commit -m tag-release-${RELEASE_TAG} k8s/base
	echo commited tag-release-${RELEASE_TAG}
elif [ ${diffCode} -eq 0 ] ; then
	echo ${RELEASE_TAG} already committed
else
	echo "git diff-index failed"
	exit 1
fi


# in any event, make sure we've tagged it locally and in-repo
git tag -a -f -m release-${RELEASE_TAG} ${RELEASE_TAG}
git push origin -f ${RELEASE_TAG}

