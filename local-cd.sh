#!/bin/bash

set -o errexit -o nounset
CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build .
TRAVIS_BRANCH=master
if [ "$TRAVIS_BRANCH" != "master" ]
then 
	    echo "This commit was made against the $TRAVIS_BRANCH and not the master! No deploy!" 
		exit 1
fi

rev=$(git rev-parse --short HEAD)
docker login -u docker-image-builder  http://hub.sky-cloud.net
docker build -t hub.sky-cloud.net/nap2/netd:${TRAVIS_BRANCH}_build-${rev} .
docker push hub.sky-cloud.net/nap2/netd:${TRAVIS_BRANCH}_build-${rev}


