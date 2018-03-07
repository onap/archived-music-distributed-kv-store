#!/bin/bash

CONSUL_IP="localhost"
MOUNTPATH="/configs/"

docker run -e CONSUL_IP=$CONSUL_IP -e MOUNTPATH=$MOUNTPATH -it --name dkv -p 8200:8200 -p 8080:8080 dkv
