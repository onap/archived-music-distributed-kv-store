#!/bin/bash

BUILD_ARGS="--no-cache"
ORG="onap"
PROJECT="music"
IMAGE="distributed-kv-store"
DOCKER_REPOSITORY="nexus3.onap.org:10003"
IMAGE_NAME="${DOCKER_REPOSITORY}/${ORG}/${PROJECT}/${IMAGE}"
TAG_NAME=${UNIQUE_DOCKER_TAG}

echo TAG_NAME: ${TAG_NAME}
set -x
if [ $HTTP_PROXY ]; then
    BUILD_ARGS+=" --build-arg HTTP_PROXY=${HTTP_PROXY}"
fi
if [ $HTTPS_PROXY ]; then
    BUILD_ARGS+=" --build-arg HTTPS_PROXY=${HTTPS_PROXY}"
fi

function generate_binary {
    pushd ../src/dkv
    make build
    popd
    cp ../target/dkv .

    # Change the following work around for reading token_service.json
    # cp ../src/dkv/api/token_service_map.json .
}

function build_image {
    echo "Start build docker image."
    docker build ${BUILD_ARGS} -t ${IMAGE_NAME}:latest .
}

function push_image {
    echo "Start push docker image."
    docker push ${IMAGE_NAME}:latest
    if [ ! -z "${TAG_NAME}" ]; then
      docker tag ${IMAGE_NAME}:latest ${IMAGE_NAME}:1.0-SNAPSHOT-${TAG_NAME}
      docker push ${IMAGE_NAME}:1.0-SNAPSHOT-${TAG_NAME}
    fi
}

function remove_binary {
    rm dkv
}

generate_binary
build_image
push_image
remove_binary
