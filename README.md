#### Rate limiting solution between services. 

This example includes:

**bucket_quoter** package:
- implementation of _[bucket quota limiter](https://github.com/alexgaas/rate_limiter/blob/main/bucket_quoter/README.md)_

**bucket_quoter_example** consists of:
- **[quoter](https://github.com/alexgaas/rate_limiter/blob/main/bucket_quoter_example/README.md)** - backend (service built with _gin-gonic_) for limiter (built also as docker image)
- **[lambda](https://github.com/alexgaas/rate_limiter/blob/main/bucket_quoter_example/lambda/README.md)** - basic AWS lambda (built with AWS Go SDK) to send requests to **quoter** (built as docker image with Kong as API Gateway) for simulation
- **[docker](https://github.com/alexgaas/rate_limiter/blob/main/bucket_quoter_example/README.md)** to run simulation of rate limiting between **quoter** and **basic_lambda**


