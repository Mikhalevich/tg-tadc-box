#! /usr/bin/env bash

TAG=$1
VERSION="${TAG#v}"

echo "deploy version ${VERSION}"

helm secrets upgrade --install do-helm-app ./script/k8s/do-helm-app --take-ownership -f ./script/k8s/do-helm-app/secrets.yaml \
--set dbmigration.image.tag=${VERSION} \
--set bot.image.tag=${VERSION} \
--set outboxpoller.image.tag=${VERSION} \
--set poller.image.tag=${VERSION}
