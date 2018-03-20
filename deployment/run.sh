#!/bin/bash

DATASTORE="consul"
DATASTORE_IP="localhost"

MOUNTPATH="/dkv_mount_path/configs/"
DEFAULT_CONFIGS=$(pwd)/../mountpath/default # TODO(sshank): Change this to think from Kubernetes Volumes perspective.

docker run -e DATASTORE=$DATASTORE -e DATASTORE_IP=$DATASTORE_IP -e MOUNTPATH=$MOUNTPATH -it \
           --name dkv \
           -v $DEFAULT_CONFIGS:/dkv_mount_path/configs/default \
           -p 8200:8200 -p 8080:8080 nexus3.onap.org:10003/onap/music/distributed-kv-store
