#! /usr/bin/env bash

sops --decrypt -i script/k8s/do-helm-app/secrets.yaml
