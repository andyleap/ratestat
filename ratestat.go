// ratestat project ratestat.go
package ratestat

import (
	"sync/atomic"
	"time"
)

type RateStat struct {
	Buckets   []uint64
	curBucket int32
	interval  time.Duration
	length    int32
}

func (rs *RateStat) Log(count uint64) {
	bucket := atomic.LoadInt32(&rs.curBucket)
	atomic.AddUint64(&rs.Buckets[bucket], count)
}

func (rs *RateStat) Value() uint64 {
	val := uint64(0)
	for k := range rs.Buckets {
		val += atomic.LoadUint64(&rs.Buckets[k])
	}
	return val
}

func (rs *RateStat) manage() {
	for range time.Tick(rs.interval) {
		bucket := atomic.LoadInt32(&rs.curBucket)
		bucket++
		if bucket >= rs.length {
			bucket = 0
		}
		atomic.StoreUint64(&rs.Buckets[bucket], 0)
		atomic.StoreInt32(&rs.curBucket, bucket)
	}
}

func New(buckets int32, interval time.Duration) *RateStat {
	rs := &RateStat{
		Buckets:  make([]uint64, buckets),
		interval: interval,
		length:   buckets,
	}
	go rs.manage()
	return rs
}
