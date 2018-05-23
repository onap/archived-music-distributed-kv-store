FROM alpine:3.7

ARG HTTP_PROXY=${HTTP_PROXY}
ARG HTTPS_PROXY=${HTTPS_PROXY}

ENV http_proxy $HTTP_PROXY
ENV https_proxy $HTTPS_PROXY
ENV DATASTORE_IP $DATASTORE_IP
ENV CONSUL_VERSION 1.0.6

RUN apk update && \
    # apk add g++ && \
    # apk add make && \
    apk add unzip && \
    # apk add git && \
    # apk add curl && \
    apk add wget && \
    apk add --no-cache bash

RUN wget -qO /tmp/consul.zip "https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip" && \
    unzip -d /bin /tmp/consul.zip && \
    chmod 755 /bin/consul && \
    rm /tmp/consul.zip

EXPOSE 8200
EXPOSE 8080


RUN mkdir /dkv_mount_path && \
    mkdir /dkv_mount_path/api && \
    mkdir /dkv_mount_path/consul_data && \
    mkdir /dkv_mount_path/configs && \
    mkdir /dkv_mount_path/configs/default

WORKDIR /dkv_mount_path
ADD ./dkv /dkv_mount_path/
# ADD ./token_service_map.json /dkv_mount_path/api/
ADD ./docker-entrypoint.sh /dkv_mount_path/

ENTRYPOINT /dkv_mount_path/docker-entrypoint.sh