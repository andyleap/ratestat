// ratestat project ratestat.go
package ratestat

import (
	"sync/atomic"
	"time"
)

type RateStat32 struct {
	Buckets   []uint32
	curValue  uint64
	curBucket int32
	interval  time.Duration
	length    int32
}

func (rs *RateStat32) Log(count uint32) {
	bucket := atomic.LoadInt32(&rs.curBucket)
	atomic.AddUint32(&rs.Buckets[bucket], count)
	atomic.AddUint64(&rs.curValue, uint64(count))
}

func (rs *RateStat32) Value() uint64 {
	return atomic.LoadUint64(&rs.curValue)
}

func (rs *RateStat32) manage() {
	for range time.Tick(rs.interval) {
		bucket := atomic.LoadInt32(&rs.curBucket)
		bucket++
		if bucket >= rs.length {
			bucket = 0
		}
		old := atomic.SwapUint32(&rs.Buckets[bucket], 0)
		atomic.StoreInt32(&rs.curBucket, bucket)
		atomic.AddUint64(&rs.curValue, -uint64(old))
	}
}

func New32(buckets int32, interval time.Duration) *RateStat32 {
	rs := &RateStat32{
		Buckets:  make([]uint32, buckets),
		interval: interval,
		length:   buckets,
	}
	go rs.manage()
	return rs
}
