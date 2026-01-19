#### Implementation of _bucket quota limiter_.

Simple example of usage:
```go
quoter := NewBucketQuoter(inflow, capacity, true, nil)
for {
    // get quota request

    if !quoter.IsAvailable() {
        // do something else:
        // quoter.Sleep()
        // or return 429
    }

    quoter.Use(1)
    // send 200
}
```

#### Pros:

- Simple (_only 250 lines of code_), production-ready and easily injectable to any Go service
- Has a stat counters that could be exposed as metrics for observations and shown on the dashboards

#### Cons:

- When call method _FillBucket_ **BucketQuoter** calculates tokens which will be added into Buckets from the time of last update. If this quantity is not
whole number that quantity will be rounded down. That may impact small **RPS**.
- When **BucketQuoter** gets time over _GetWaitTime_, it returns not the nearest time when new tokens available but time proportional to 1 / **RPS** in the microseconds.
Since rounding down may happen due division operation, _GetWaitTime_ may return a time when tokens not accrued yet.


