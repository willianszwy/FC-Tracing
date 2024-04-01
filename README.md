# FC-Tracing

In the directory of ServiceB copy the .env.example to .env and change the values

## Build
```shell
docker-compose build
```

## Run
```shell
docker-compose up
```

## request the endpoint POST

```shell
curl -X POST http://localhost:8081 -d '{"zipcode": "06835100"}'
```


## Zipkin
http://127.0.0.1:9411/zipkin/