## simple redis client

```
# REDIS

$ docker run --rm -d --name redis -p 6379:6379 redis


# CLIENT

$ make go-run


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