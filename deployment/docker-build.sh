#!/bin/bash

BUILD_ARGS="--no-cache"
ORG="onap"
VERSION="1.0.0"
PROJECT="music"
IMAGE="distributed-kv-store"
DOCKER_REPOSITORY="nexus3.onap.org:10003"
IMAGE_NAME="${DOCKER_REPOSITORY}/${ORG}/${PROJECT}/${IMAGE}"
TIMESTAMP=$(date +"%Y%m%dT%H%M%S")

if [ $HTTP_PROXY ]; then
    BUILD_ARGS+=" --build-arg HTTP_PROXY=${HTTP_PROXY}"
fi
if [ $HTTPS_PROXY ]; then
    BUILD_ARGS+=" --build-arg HTTPS_PROXY=${HTTPS_PROXY}"
fi

function install_build_dependencies {
    local golang_version=go1.10.linux-amd64
    sudo apt-get install -y build-essential
    if [ ! -d /opt/go ]; then
        mkdir /opt/go
        pushd /opt/go
        curl -O https://dl.google.com/go/$golang_version.tar.gz
        tar -zxf $golang_version.tar.gz
        echo GOROOT=$PWD/go >> /etc/environment
        echo PATH=$PATH:$PWD/go/bin >> /etc/environment
        rm -rf tar -zxf $golang_version.tar.gz
        popd
    fi
    source /etc/environment
}

function generate_binary {
    pushd ../src/dkv
    make build
    popd
    cp ../target/dkv .
}

function build_image {
    echo "Start build docker image."
    docker build ${BUILD_ARGS} -t ${IMAGE_NAME}:latest .
}

function push_image {
    echo "Start push docker image."
    docker push ${IMAGE_NAME}:latest
}

function remove_binary {
    rm dkv
}

install_build_dependencies
generate_binary
build_image
push_image
remove_binary
