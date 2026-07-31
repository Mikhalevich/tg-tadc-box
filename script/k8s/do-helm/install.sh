#! /usr/bin/env bash

TAG=$1
VERSION="${TAG#v}"

echo "deploy version ${VERSION}"

helm upgrade --install jaeger ./script/k8s/do-helm/jaeger
helm secrets upgrade --install bot ./script/k8s/do-helm/bot -f ./script/k8s/do-helm/bot/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}
helm secrets upgrade --install poller ./script/k8s/do-helm/poller -f ./script/k8s/do-helm/poller/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}
helm secrets upgrade --install outboxpoller ./script/k8s/do-helm/outboxpoller -f ./script/k8s/do-helm/outboxpoller/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}
