package models

import (
	"database/sql"
	"sync/atomic"
	"time"
)

type DBMetrics struct {
	OpenConnections   int32
	InUseConnections int32
	WaitCount        int64
	WaitDuration     time.Duration
	MaxIdleTimeClosed int64
}

type MetricsCollector struct {
	db      *sql.DB
	metrics DBMetrics
}

func NewMetricsCollector(db *sql.DB) *MetricsCollector {
	mc := &MetricsCollector{db: db}
	go mc.collect()
	return mc
}

func (mc *MetricsCollector) collect() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		stats := mc.db.Stats()
		
		atomic.StoreInt32(&mc.metrics.OpenConnections, int32(stats.OpenConnections))
		atomic.StoreInt32(&mc.metrics.InUseConnections, int32(stats.InUse))
		atomic.StoreInt64(&mc.metrics.WaitCount, stats.WaitCount)
		atomic.StoreInt64(&mc.metrics.MaxIdleTimeClosed, stats.MaxIdleClosed)
	}
}

func (mc *MetricsCollector) GetMetrics() DBMetrics {
	return DBMetrics{
		OpenConnections:   atomic.LoadInt32(&mc.metrics.OpenConnections),
		InUseConnections: atomic.LoadInt32(&mc.metrics.InUseConnections),
		WaitCount:        atomic.LoadInt64(&mc.metrics.WaitCount),
		MaxIdleTimeClosed: atomic.LoadInt64(&mc.metrics.MaxIdleTimeClosed),
	}
}