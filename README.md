#### Rate limiting solution between services. 

This example includes:

**bucket_quoter** package:
- implementation of _[bucket quota limiter](https://github.com/alexgaas/rate_limiter)_

**bucket_quoter_example** consists of:
- **[quoter]((https://github.com))** - backend (service built with _gin-gonic_) for limiter (built also as docker image).
- **[lambda](https://github.com)** - basic AWS lambda (built with AWS Go SDK) to send requests to **quoter** (built as docker image with Kong as API Gateway) for simulation
- **[docker](https://github.com)** to run simulation. You can see [results](https://github.com/alexgaas/rate_limiter) of experiment of rate limiting between **quoter** and **basic_lambda**.


