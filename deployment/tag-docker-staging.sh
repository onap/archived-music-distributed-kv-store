#!/bin/bash

ORG="onap"
PROJECT="music"
IMAGE="distributed-kv-store"
DOCKER_REPOSITORY="nexus3.onap.org:10003"
IMAGE_NAME="${DOCKER_REPOSITORY}/${ORG}/${PROJECT}/${IMAGE}"
TAG_NAME=${UNIQUE_DOCKER_TAG}

echo TAG_NAME: ${TAG_NAME}
set -x

echo "Push STAGING tag for docker image."
if [ ! -z "${TAG_NAME}" ]; then
    docker pull ${IMAGE_NAME}:1.0-SNAPSHOT-${TAG_NAME}
    docker tag ${IMAGE_NAME}:1.0-SNAPSHOT-${TAG_NAME} ${IMAGE_NAME}:1.0-STAGING-latest
    docker push ${IMAGE_NAME}:1.0-STAGING-latest
fi

