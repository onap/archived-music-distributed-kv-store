#!/bin/bash

function start_consul_server {
    consul agent -bootstrap -server -bind=127.0.0.1 -data-dir=/dkv_mount_path/consul_data &
}

function start_api_server {
    pushd /dkv_mount_path/
    ./dkv
}

if [ "$DATASTORE_IP" = "localhost" ]; then
    start_consul_server
    sleep 5
fi
start_api_server
