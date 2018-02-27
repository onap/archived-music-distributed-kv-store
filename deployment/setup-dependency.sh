#!/bin/bash

function install_go {
    local golang_version=go1.10.linux-amd64
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

function install_dependencies {
    pushd src/dkv/
    make all
    popd
}

function create_mountpath {
    cp -r mountpath/ /configs
}

install_go
install_dependencies
create_mountpath
