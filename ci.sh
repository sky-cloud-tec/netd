#!/bin/bash

set -o errexit -o nounset
TRAVIS_BRANCH=master
if [ "$TRAVIS_BRANCH" != "feat_config_based_engine" ]
then 
	    echo "This commit was made against the $TRAVIS_BRANCH and not the master! No deploy!" 
		exit 1
fi

rev=$(git rev-parse --short HEAD)
num=$(git rev-list --count feat_config_based_engine)
name=registry.cn-hangzhou.aliyuncs.com/sky-cloud-tec/netd:${TRAVIS_BRANCH}-travis-${rev}-${num}
sudo cat /etc/docker/daemon.json | jq '. + {"insecure-registries": ["registry.cn-hangzhou.aliyuncs.com"]}' | sudo tee /etc/docker/daemon.json
sudo service docker restart
echo "$HUB" | docker login -u $USER  registry.cn-hangzhou.aliyuncs.com --password-stdin
docker build -t ${name} .
docker push ${name}
