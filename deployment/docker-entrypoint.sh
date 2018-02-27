#!/bin/bash

function install_go {
    mkdir /go
    pushd /go
    curl -O https://dl.google.com/go/go1.10.linux-amd64.tar.gz
    tar -zxf go1.10.linux-amd64.tar.gz
    rm go1.10.linux-amd64.tar.gz
    export GOROOT=$PWD/go
    export PATH=$PATH:$GOROOT/bin
    popd
}

function install_dependencies {
    pushd src/dkv/
    make all
    popd
}

function verify_consul_run {
    consul --version
}

function set_go_path {
    export GOPATH=$PWD
}

function start_consul_server {
    # Running consul in server mode since we are doing a single node. If we need to add more,
    # We need to run multiple consul agents in client mode without providing the -server arguements.

    # CHANGE THIS TO SERVER MODE!
    consul agent -dev > /dev/null 2>&1 &
    # consul agent -server -bind=127.0.0.1 -data-dir=/dkv/consul &
}

function start_api_server {
    go run src/dkv/main.go
}

install_go
install_dependencies
set_go_path
start_consul_server
start_api_server
