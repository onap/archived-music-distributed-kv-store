#!/bin/bash

CONSUL_IP=localhost

docker run -e CONSUL_IP=$CONSUL_IP -it --name dkv -p 8200:8200 -p 8080:8080 dkv
