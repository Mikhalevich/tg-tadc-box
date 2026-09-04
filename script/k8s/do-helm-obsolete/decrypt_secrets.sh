#! /usr/bin/env bash

sops --decrypt -i script/k8s/do-helm/bot/secrets.yaml
sops --decrypt -i script/k8s/do-helm/poller/secrets.yaml
sops --decrypt -i script/k8s/do-helm/outboxpoller/secrets.yaml
