Distributed Key Value Store using Consul to store application configuration data.

# TODO (Add documentation on how to run)

# Sample Curl examples:

## Register project
`curl -X POST -d '{"domain":"<project>"}' localhost:8080/v1/register`
`export TOKEN=`
## Register sub project 
`curl -X POST -d '{"subdomain":"<sub-project>"}' localhost:8080/v1/register/$TOKEN/subdomain`

## Upload properties file
`curl -X POST -F 'token=$TOKEN' -F 'configFile=@./example.properties' localhost:8080/v1/config`
`curl -X POST -F 'token=$TOKEN' -F 'subdomain=<sub-project>' -F 'configFile=@./example.properties' localhost:8080/v1/config`

## Load properties file into Consul
`curl -X POST -d '{"token":"$TOKEN", "filename": "example.properties"}' localhost:8080/v1/config/load`

## Fetch properties file
`curl -X GET localhost:8080/v1/config/$TOKEN/example.properties`
`curl -X GET localhost:8080/v1/config/$TOKEN/<sub-project>/example.properties`

## Delete properties file
`curl -X DELETE localhost:8080/v1/config/$TOKEN/example.properties`
`curl -X DELETE localhost:8080/v1/config/$TOKEN/<sub-project>/example.properties`

## Delete project/sub project
`curl -X DELETE localhost:8080/v1/register/$TOKEN/subdomain/<sub-project>`
`curl -X DELETE localhost:8080/v1/register/$TOKEN`
