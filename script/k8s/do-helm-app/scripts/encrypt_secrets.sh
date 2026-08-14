#! /usr/bin/env bash

sops --encrypt -i script/k8s/do-helm-app/secrets.yaml
