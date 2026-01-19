#### Docker configuration for experiment:

- Basic **AWS** lambda image
- **Kong** to run a Docker based Lambda function and access it locally via Kong API Gateway to test it, 
the same way one would normally do if the function was behind an AWS API Gateway. 
Please see this [githib repo](https://github.com/brafales/docker-lambda-kong) with details.
- Quoter image

#### Experiment:

Prerequisites:
- You need to have **Docker** installed
- Build docker images for both **quoter** service and basic **lambda**

Run simulation:

- Run docker with compose as `docker compose up -d`
- _quoter_ service have default setup for a bar **0.1 RPS** (5 requests per minute) for corresponding **subscription ID** 
- call
```go
curl localhost:8000 -d '{}'
```
few times. First 5 times you have to get back **200 OK** but next **429** status code back (`Too Many Requests`)