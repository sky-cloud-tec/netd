#!/bin/bash

set -o errexit -o nounset
TRAVIS_BRANCH=master
if [ "$TRAVIS_BRANCH" != "master" ]
then 
	    echo "This commit was made against the $TRAVIS_BRANCH and not the master! No deploy!" 
		exit 1
fi

rev=$(git rev-parse --short HEAD)
num=$(git rev-list --count master)
name=hub.sky-cloud.net/nap2/netd:${TRAVIS_BRANCH}_build-${rev}-${num}
sudo cat /etc/docker/daemon.json | jq '. + {"insecure-registries": ["hub.sky-cloud.net"]}' | sudo tee /etc/docker/daemon.json
sudo service docker restart
echo "$HUB" | docker login -u docker-image-builder  http://hub.sky-cloud.net --password-stdin
docker build -t ${name} .
docker push ${name}
