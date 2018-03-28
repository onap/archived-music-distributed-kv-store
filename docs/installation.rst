.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

Installation
============

As for current release, project is packaged as a single Docker Container. For subsequent
releases, it will be integrated with OOM.

.. code-block:: console

    # Set Datastore as Consul
    DATASTORE="consul"
    # Set IP address of where Consul is running
    DATASTORE_IP="localhost"
    # Set mountpath inside the container where persistent data is stored.
    MOUNTPATH="/dkv_mount_path/configs/"
    # Place all Config data which needs to be loaded in default directory.
    DEFAULT_CONFIGS=$(pwd)/mountpath/default
    # Create the directories.
    mkdir -p mountpath/default
    # Login to Nexus.
    docker login -u docker -p docker nexus3.onap.org:10001
    # Pull distributed-kv-store image.
    docker pull nexus3.onap.org:10001/onap/music/distributed-kv-store
    # Run the distributed-kv-store image.
    docker run -e DATASTORE=$DATASTORE -e DATASTORE_IP=$DATASTORE_IP -e MOUNTPATH=$MOUNTPATH -d \
           --name dkv \
           -v $DEFAULT_CONFIGS:/dkv_mount_path/configs/default \
           -p 8200:8200 -p 8080:8080 nexus3.onap.org:10001/onap/music/distributed-kv-store

.. end
