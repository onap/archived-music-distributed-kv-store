#!/bin/bash
set -x
echo TEST SCRIPT WITH PARAMETER ${1}
export CSIT_TEST_VAR="has value"
export UNIQUE_DOCKER_TAG="WAS CHANGED BY THE SCRIPT"
