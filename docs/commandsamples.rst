.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

Sample Commands
===============

.. code-block:: console

    ## Load default configuration
    curl -X GET localhost:8080/v1/config/load-default

    ## Check if Keys were loaded into Consul
    curl -X GET localhost:8080/v1/getconfigs

    ## Check value for a single key
    curl -X GET localhost:8080/v1/getconfig/<key>

    ## Register new domain
    curl -X POST -d '{"domain":"new_project"}' localhost:8080/v1/register
    export TOKEN=
    ## Register new sub domain
    curl -X POST -d '{"subdomain":"sub_project"}' localhost:8080/v1/register/$TOKEN/subdomain

    ## Check if a domain is already registered.
    curl -X GET localhost:8080/v1/register/$TOKEN

    ## Upload properties file to domain or subdomain.
    curl -X POST -F 'token=$TOKEN' -F 'configFile=@./example.properties' localhost:8080/v1/config
    curl -X POST -F 'token=$TOKEN' -F 'subdomain=sub_domain' -F 'configFile=@./example.properties' localhost:8080/v1/config

    ## Load properties file into Consul
    curl -X POST -d '{"token":"$TOKEN", "filename": "example.properties"}' localhost:8080/v1/config/load

    ## Fetch properties file
    curl -X GET localhost:8080/v1/config/$TOKEN/example.properties
    curl -X GET localhost:8080/v1/config/$TOKEN/sub_domain/example.properties

    ## Delete properties file
    curl -X DELETE localhost:8080/v1/config/$TOKEN/example.properties
    curl -X DELETE localhost:8080/v1/config/$TOKEN/sub_domain/example.properties

    ## Delete project/sub project
    curl -X DELETE localhost:8080/v1/register/$TOKEN/sub_domain/sub-domain
    curl -X DELETE localhost:8080/v1/register/$TOKEN

.. end
