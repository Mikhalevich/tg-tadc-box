#! /usr/bin/env bash

sops --encrypt -i script/k8s/minikube-helm/bot/secrets.yaml
sops --encrypt -i script/k8s/minikube-helm/outboxpoller/secrets.yaml
sops --encrypt -i script/k8s/minikube-helm/poller/secrets.yaml
sops --encrypt -i script/k8s/minikube-helm/postgres/secrets.yaml
