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