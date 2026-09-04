#! /usr/bin/env bash

TAG=$1
VERSION="${TAG#v}"

echo "load image with version ${VERSION}"

docker build -t bot:${VERSION} -f ./script/docker/bot.Dockerfile .
docker tag bot:${VERSION} registry.digitalocean.com/tadc-box-registry/bot:${VERSION}
docker push registry.digitalocean.com/tadc-box-registry/bot:${VERSION}

docker build -t sqlmigrate:${VERSION} -f ./script/docker/sqlmigrate.Dockerfile .
docker tag sqlmigrate:${VERSION} registry.digitalocean.com/tadc-box-registry/sqlmigrate:${VERSION}
docker push registry.digitalocean.com/tadc-box-registry/sqlmigrate:${VERSION}

docker build -t poller:${VERSION} -f ./script/docker/poller.Dockerfile .
docker tag poller:${VERSION} registry.digitalocean.com/tadc-box-registry/poller:${VERSION}
docker push registry.digitalocean.com/tadc-box-registry/poller:${VERSION}

docker build -t outboxpoller:${VERSION} -f ./script/docker/outboxpoller.Dockerfile .
docker tag outboxpoller:${VERSION} registry.digitalocean.com/tadc-box-registry/outboxpoller:${VERSION}
docker push registry.digitalocean.com/tadc-box-registry/outboxpoller:${VERSION}
