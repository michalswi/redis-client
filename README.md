## simple redis client

Simple redis client written in Go. You can interact (set, get) either with local redis running in docker or the Azure Cache for Redis.

#### # image

`michalsw/redis-client:latest`

While connecting to Redis by default there is no password set in **redis-client**. If you want to pass access key (password) use env **REDIS_PASS**. If TLS is required like in Azure use env **REDIS_TLS**.  

Use **make** to make you better.

```
$ make

Usage:
  make <target>

Targets:
  go-run           Run redis client - no binary
  go-build         Build binary
  docker-build     Build docker image
  docker-run       Once docker image is ready run with default parameters (or overwrite)
  docker-stop      Stop running docker
  azure-rg         Create the Azure Resource Group
  azure-rg-del     Delete the Azure Resource Group
  azure-aci        Run redis-client app (Azure Container Instance)
  azure-aci-logs   Get redis-client app logs (Azure Container Instance)
  azure-aci-delete  Delete redis-client app (Azure Container Instance)
```

#### # endpoints
```
GET   /red/ping
POST  /red/setuser
GET   /red/getuser/{id}
GET   /red/home
```

#### # local
```
# deploy 'redis'

$ docker run --rm -d --name redis -p 6379:6379 redis


# run 'redis-client'

$ make go-run

OR

$ REDIS_HOST=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis)
$ REDIS_HOST=$REDIS_HOST make docker-run


# POST

$ curl -i -X POST -d '{"name":"mi","age":10}' localhost:8080/red/setuser
$ curl -i -X POST -d '{"name":"mo","age":20}' localhost:8080/red/setuser


# GET

$ curl -XGET localhost:8080/red/ping
PONG

$ curl -XGET localhost:8080/red/getuser/1 | jq
{
  "name": "mi",
  "age": 10
}

$ curl -XGET localhost:8080/red/getuser/2 | jq
{
  "name": "mo",
  "age": 20
}
```

#### # azure

[Here](https://docs.microsoft.com/en-us/azure/azure-cache-for-redis/cache-python-get-started) you have some example how to create a Python app that uses Azure Cache for Redis.

```
### deploy 'redis'

ACI - Azure Container Instances

## ACI - NO access keys (password)

# > redis

$ make azure-rg

$ DNS_NAME_LABEL=redis-$RANDOM \
  LOCATION=westeurope \
  RGNAME=redisRG

$ az container create \
  --resource-group $RGNAME \
  --name redis \
  --image redis \
  --restart-policy Always \
  --ports 6379 \
  --dns-name-label $DNS_NAME_LABEL \
  --location $LOCATION \
  --environment-variables \
    DNS_NAME=$DNS_NAME_LABEL.$LOCATION.azurecontainer.io


# > test using 'redis-client' (run locally)

$ REDIS_HOST=redis-28318.westeurope.azurecontainer.io make go-run

$ curl -XGET localhost:8080/red/ping


## TERRAFORM - WITH access keys (password)

# > redis
in progress..


# > test using redis-client (run locally)

$ REDIS_HOST=redisms.redis.cache.windows.net \
REDIS_PORT=6380 \
REDIS_PASS=yIvSS+rPz3zWhG3685lj6Fw9Si51stlZgx4lYieWF0s= \
REDIS_TLS=true \
make go-run

$ curl -i -XGET localhost:8080/red/ping


# > test using redis-client using ACI

$ SERVICE_ADDR=80 \
REDIS_HOST=redisms.redis.cache.windows.net \
REDIS_PORT=6380 \
REDIS_PASS=yIvSS+rPz3zWhG3685lj6Fw9Si51stlZgx4lYieWF0s= \
REDIS_TLS=true \
make azure-aci

$ curl redis-client-d58df48.westeurope.azurecontainer.io/red/home

$ curl redis-client-d58df48.westeurope.azurecontainer.io/red/ping

$ make azure-aci-logs 

$ make azure-aci-delete
```