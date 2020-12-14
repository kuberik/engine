#!/bin/sh

export HOME=/kaniko
# "the command stores credentials in $HOME/.docker/config.json"
echo ${DOCKER_REGISTRY_PASSWORD} | docker login --username ${DOCKER_REGISTRY_USERNAME} --password-stdin ${DOCKER_REGISTRY_SERVER}
