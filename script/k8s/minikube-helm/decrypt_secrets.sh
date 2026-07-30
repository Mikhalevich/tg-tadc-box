#! /usr/bin/env bash

sops --decrypt -i script/k8s/minikube-helm/bot/secrets.yaml
sops --decrypt -i script/k8s/minikube-helm/outboxpoller/secrets.yaml
sops --decrypt -i script/k8s/minikube-helm/poller/secrets.yaml
sops --decrypt -i script/k8s/minikube-helm/postgres/secrets.yaml
