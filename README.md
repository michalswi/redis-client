## simple redis client

```
# REDIS

$ docker run --rm -d --name redis -p 6379:6379 redis


# RUN

$ make go-run

OR

$ REDIS_HOST=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis)
$ REDIS_HOST=$REDIS_HOST make docker-run


# POST

$ curl -i -X POST -d '{"name":"mi","age":10}' localhost:8080/red/setuser
$ curl -i -X POST -d '{"name":"mo","age":20}' localhost:8080/red/setuser


# GET

$ curl localhost:8080/red/getuser/1 | jq
{
  "name": "mi",
  "age": 10
}

$ curl localhost:8080/red/getuser/2 | jq
{
  "name": "mo",
  "age": 20
}
```