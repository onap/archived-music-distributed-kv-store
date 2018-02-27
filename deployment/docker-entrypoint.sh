#!/bin/bash

function verify_consul_run {
    consul --version
}

function start_consul_server {
    # Running consul in server mode since we are doing a single node. If we need to add more,
    # We need to run multiple consul agents in client mode without providing the -server arguements.

    # CHANGE THIS TO SERVER MODE!
    # consul agent -dev > /dev/null 2>&1 &
    consul agent -bootstrap -server -bind=127.0.0.1 -data-dir=/dkv/consul &
}

function start_api_server {
    # Uncomment the following after the mountpath is setup in the code base and the docker file.
    # Until then, go run is used.
    #cd target
    #./dkv
    cd src/dkv/
    go run main.go
}

function set_paths {
    export GOPATH=$PWD
    source /etc/environment
}

set_paths
start_consul_server
sleep 5
start_api_server
