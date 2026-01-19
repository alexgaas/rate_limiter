package bucket_quoter

import (
	"sync"
	"sync/atomic"
	"time"
)

// Token Bucket

type BucketQuoter struct {
	bucketMutex sync.Mutex
	timer       InstantTimer

	Bucket int64
	SeqNo  int64

	LastAdd int64

	FixedInflow   atomic.Int64
	FixedCapacity atomic.Int64

	InflowTokensPerSecond *atomic.Int64
	BucketTokensCapacity  *atomic.Int64

	Stat *BucketQuoterStat
}

type Result struct {
	Before int64
	After  int64
	SeqNo  int64
}

func NewBucketQuoter(inflow int64, capacity int64, fill bool, stat *BucketQuoterStat) *BucketQuoter {
	timer := NewInstantTimerMs()
	if stat == nil {
		stat = &BucketQuoterStat{}
	}

	var bucket int64 = 0
	if fill {
		bucket = atomic.LoadInt64(&capacity)
	}

	var fixedCapacity atomic.Int64
	fixedCapacity.Add(capacity)

	var fixedInflow atomic.Int64
	fixedInflow.Add(inflow)

	return &BucketQuoter{
		bucketMutex:           sync.Mutex{},
		timer:                 timer,
		Bucket:                bucket,
		SeqNo:                 0,
		LastAdd:               timer.Now(),
		FixedInflow:           fixedInflow,
		FixedCapacity:         fixedCapacity,
		InflowTokensPerSecond: &fixedInflow,
		BucketTokensCapacity:  &fixedCapacity,
		Stat:                  stat,
	}
}

// PUBLIC

func (q *BucketQuoter) IsAvailable() bool {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	q.fillBucket()
	if q.Bucket < 0 {
		// stat
		q.Stat.BucketUnderflows += 1
	}

	return q.Bucket >= 0
}

func (q *BucketQuoter) IsAvailableWithResult(r *Result) bool {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	r.Before = q.Bucket
	q.fillBucket()
	if q.Bucket < 0 {
		// stat
		q.Stat.BucketUnderflows += 1
	}
	r.After = q.Bucket
	r.SeqNo = q.SeqNo + 1

	return q.Bucket >= 0
}

func (q *BucketQuoter) GetAvailable() int64 {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	q.fillBucket()

	if q.Bucket > 0 {
		return q.Bucket
	}
	return 0
}

func (q *BucketQuoter) GetAvailableWithResult(r *Result) int64 {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	r.Before = q.Bucket
	q.fillBucket()
	r.After = q.Bucket
	r.SeqNo = q.SeqNo + 1

	if q.Bucket > 0 {
		return q.Bucket
	}
	return 0
}

func (q *BucketQuoter) Use(tokens int64) {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	q.useNoLock(tokens, false)
}

func (q *BucketQuoter) UseWithSleep(tokens int64) {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	q.useNoLock(tokens, true)
}

func (q *BucketQuoter) UseWithResult(tokens int64, r *Result, sleep bool) {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	r.Before = q.Bucket
	q.useNoLock(tokens, sleep)
	r.After = q.Bucket
	r.SeqNo = q.SeqNo + 1
}

func (q *BucketQuoter) UseAndFill(tokens int64) int64 {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	q.useNoLock(tokens, false)
	q.fillBucket()

	return q.Bucket
}

func (q *BucketQuoter) Add(tokens int64) {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	q.addNoLock(tokens)
}

func (q *BucketQuoter) AddWithResult(tokens int64, r *Result) {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	r.Before = q.Bucket
	q.addNoLock(tokens)
	r.After = q.Bucket
	r.SeqNo = q.SeqNo + 1
}

func (q *BucketQuoter) GetWaitTime() int64 {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	q.fillBucket()
	if q.Bucket >= 0 {
		return 0
	}

	return (-q.Bucket * 1000000) / q.InflowTokensPerSecond.Load()
}

func (q *BucketQuoter) GetWaitTimeWithResult(r *Result) int64 {
	q.bucketMutex.Lock()
	defer q.bucketMutex.Unlock()

	r.Before = q.Bucket
	q.fillBucket()
	r.After = q.Bucket
	r.SeqNo = q.SeqNo + 1

	if q.Bucket >= 0 {
		return 0
	}

	return (-q.Bucket * 1000000) / q.InflowTokensPerSecond.Load()
}

func (q *BucketQuoter) Sleep() {
	for !q.isAvailableNoLock() {
		delay := q.GetWaitTime()
		if delay != 0 {
			time.Sleep(time.Duration(delay) * time.Microsecond)

			// stat
			q.Stat.UsecWaited += delay
		}
	}
}

// PRIVATE

func (q *BucketQuoter) isAvailableNoLock() bool {
	q.fillBucket()

	return q.Bucket >= 0
}

func (q *BucketQuoter) fillBucket() {
	timerNow := q.timer.Now()
	elapsed := q.timer.Duration(q.LastAdd, timerNow)

	if q.InflowTokensPerSecond.Load()*elapsed > q.timer.Resolution() {
		inflow := q.InflowTokensPerSecond.Load() * elapsed / q.timer.Resolution()
		if q.Stat != nil {
			q.Stat.AggregateInflow += inflow
		}

		q.Bucket += inflow

		if q.Bucket > q.BucketTokensCapacity.Load() {
			q.Bucket = q.BucketTokensCapacity.Load()
		}

		q.LastAdd = timerNow
	}
}

func (q *BucketQuoter) useNoLock(tokens int64, sleep bool) {
	if sleep {
		q.Sleep()
	}
	q.Bucket -= tokens

	// stat
	q.Stat.TokensUsed += tokens
	q.Stat.MsgPassed += 1
}

func (q *BucketQuoter) addNoLock(tokens int64) {
	q.Bucket += tokens
	if q.Bucket > q.BucketTokensCapacity.Load() {
		q.Bucket = q.BucketTokensCapacity.Load()
	}
}
