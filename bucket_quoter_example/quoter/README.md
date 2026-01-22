#### Build
In order to test that this all works as expected, we need to build that Docker image and run it:
```shell
go mod vendor

docker build -t docker-image:quoter --platform linux/arm64 .
docker run -p 9443:8443 docker-image:quoter
```

#### Local Test
Test simply calling _ping_ on running docker service
```shell
curl -k "https://localhost:9443/ping"
```
or
```shell
curl -H "X-Limiter-Subscription-ID: 897d9f58-6b42-4ca7-8229-2e04056490b7" -k "https://localhost:9443/limit"
```