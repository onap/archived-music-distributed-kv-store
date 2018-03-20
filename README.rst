Distributed Key Value Store using Consul to store application configuration data.

# TODO (Add documentation on how to run)

# Sample Curl examples:

## Load default configuration
`curl -X GET localhost:8080/v1/config/load-default`

## Register new domain
`curl -X POST -d '{"domain":"<project>"}' localhost:8080/v1/register`
`export TOKEN=`
## Register new sub domain 
`curl -X POST -d '{"subdomain":"<sub-project>"}' localhost:8080/v1/register/$TOKEN/subdomain`

## Check if a domain is already registered.
`curl -X GET localhost:8080/v1/register/$TOKEN`

## List all sub domains in a domain.
`TODO`

## Upload properties file to domain or subdomain.
`curl -X POST -F 'token=$TOKEN' -F 'configFile=@./example.properties' localhost:8080/v1/config`
`curl -X POST -F 'token=$TOKEN' -F 'subdomain=<sub-domain>' -F 'configFile=@./example.properties' localhost:8080/v1/config`

## Load properties file into Consul
`curl -X POST -d '{"token":"$TOKEN", "filename": "example.properties"}' localhost:8080/v1/config/load`

## Fetch properties file
`curl -X GET localhost:8080/v1/config/$TOKEN/example.properties`
`curl -X GET localhost:8080/v1/config/$TOKEN/<sub-domain>/example.properties`

## Delete properties file
`curl -X DELETE localhost:8080/v1/config/$TOKEN/example.properties`
`curl -X DELETE localhost:8080/v1/config/$TOKEN/<sub-domain>/example.properties`

## Delete project/sub project
`curl -X DELETE localhost:8080/v1/register/$TOKEN/subdomain/<sub-domain>`
`curl -X DELETE localhost:8080/v1/register/$TOKEN`
