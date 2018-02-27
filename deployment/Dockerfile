FROM ubuntu:14.04

ARG HTTP_PROXY=${HTTP_PROXY}
ARG HTTPS_PROXY=${HTTPS_PROXY}

ENV http_proxy $HTTP_PROXY
ENV https_proxy $HTTPS_PROXY
ENV CONSUL_IP $CONSUL_IP
ENV CONSUL_VERSION 1.0.6

# Run Docker build from dkv directory.
WORKDIR /distributed-kv-store

RUN apt-get update && \
    apt-get install -y build-essential && \
    apt-get install -y realpath && \
    apt-get install -y unzip && \
    apt-get install -y git && \
    apt-get install -y curl && \
    apt-get install -y wget && \
    git clone https://git.onap.org/music/distributed-kv-store


RUN wget -qO /tmp/consul.zip "https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip" && \
    unzip -d /bin /tmp/consul.zip && \
    chmod 755 /bin/consul && \
    rm /tmp/consul.zip

EXPOSE 8200
EXPOSE 8080

# Change this when deployment gets merged.
WORKDIR /distributed-kv-store/distributed-kv-store/deployment/
ADD ./setup-dependency.sh /distributed-kv-store/distributed-kv-store/deployment/
ADD ./docker-entrypoint.sh /distributed-kv-store/distributed-kv-store/deployment/

WORKDIR /distributed-kv-store/distributed-kv-store
RUN deployment/setup-dependency.sh

VOLUME /configs

ENTRYPOINT deployment/docker-entrypoint.sh
#ENTRYPOINT /bin/bash