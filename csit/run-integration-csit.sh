#!/bin/bash -x
#
# Copyright 2020 © Samsung Electronics Co., Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# $1 unique docker tag ("latest" by default)

set -x
export UNIQUE_DOCKER_TAG=${1:-latest}
export GERRIT_BRANCH=${GERRIT_BRANCH:-master}
export WORKSPACE=$(git rev-parse --show-toplevel)/csit
rm -rf ${WORKSPACE}/data

mkdir -p ${WORKSPACE}/data
cd ${WORKSPACE}/data

git clone https://gerrit.onap.org/r/integration/csit -b ${GERRIT_BRANCH}
cp csit/*.sh ${WORKSPACE}/
cd ${WORKSPACE}
./run-csit.sh plans/music-distributed-kv-store-test-plan
