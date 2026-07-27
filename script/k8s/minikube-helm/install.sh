#! /usr/bin/env bash

TAG=$1
VERSION="${TAG#v}"

echo "deploy version ${VERSION}"

helm secrets upgrade --install --kube-context minikube postgres ./script/k8s/minikube-helm/postgres -f ./script/k8s/minikube-helm/postgres/secrets.yaml
helm upgrade --install --kube-context minikube jaeger ./script/k8s/minikube-helm/jaeger
helm secrets upgrade --install --kube-context minikube bot ./script/k8s/minikube-helm/bot -f ./script/k8s/minikube-helm/bot/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}
helm secrets upgrade --install --kube-context minikube poller ./script/k8s/minikube-helm/poller -f ./script/k8s/minikube-helm/poller/secrets.yaml --set dbmigration.image.tag=${VERSION} --set image.tag=${VERSION}
