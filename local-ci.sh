#!/bin/bash

set -o errexit -o nounset
CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build .
TRAVIS_BRANCH=feat_config_based_engine
if [ "$TRAVIS_BRANCH" != "feat_config_based_engine" ]
then 
	    echo "This commit was made against the $TRAVIS_BRANCH and not the master! No deploy!" 
		exit 1
fi

rev=$(git rev-parse --short HEAD)
num=$(git rev-list --count feat_config_based_engine)
name=hub.sky-cloud.net/nap2/netd:${TRAVIS_BRANCH}_local-${rev}-${num}
docker login -u docker-image-builder  http://hub.sky-cloud.net
docker build -t ${name} .
docker push ${name}


