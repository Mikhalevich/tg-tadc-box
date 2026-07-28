#! /usr/bin/env bash

sops --encrypt -i script/k8s/do-helm/bot/secrets.yaml
sops --encrypt -i script/k8s/do-helm/poller/secrets.yaml
