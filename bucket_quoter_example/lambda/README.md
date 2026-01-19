#### Build
In order to test that this all works as expected, we need to build that Docker image and run it:
```shell
docker build -t docker-image:basic-lambda --platform linux/arm64 .
docker run -e DRY_RUN="1" -p 9000:8080 docker-image:basic-lambda
```

#### Local Test
The AWS provided base Docker images already come with something called the Runtime Interface Client which takes care of acting as that proxy for you, 
allowing the invocation of the function via an HTTP API call.

In order to get our local Lambda to reply with a response, this is what we need to do:
```shell
curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
```


