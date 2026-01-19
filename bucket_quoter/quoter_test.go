package bucket_quoter

import (
	"fmt"
	"testing"
	"time"
)

// BASIC TESTS

var (
	// with up to 60 seconds quota buildup.
	inflow, capacity int64 = 5, 60000
)

func TestBasicBucketQuoter(t *testing.T) {
	fmt.Printf("initial inflow %d, capacity %d\n", inflow, capacity)
	quoter := NewBucketQuoter(inflow, capacity, true, nil)
	for {
		// get message

		if !quoter.IsAvailable() {
			fmt.Printf("no more tokens available, bucket size is %d\n", quoter.Bucket)
			// do something else:

			// quoter.Sleep()
			// or
			return
		}

		quoter.Use(1)
		fmt.Printf("current bucket size is %d\n", quoter.Bucket)

		// send message
	}
}

func TestBasicBucketQuoterWithSleep(t *testing.T) {
	fmt.Printf("initial inflow %d, capacity %d\n", inflow, capacity)

	quoter := NewBucketQuoter(inflow, capacity, true, nil)
	for {
		// get message

		quoter.Sleep()
		quoter.Use(1)
		fmt.Printf("current bucket size is %d\n", quoter.Bucket)

		time.Sleep(time.Duration(100) * time.Millisecond)

		// send message
	}
}
