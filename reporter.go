// Package plugin provides provides the ability for registering plugins.
package plugin

import (
	"time"

	"github.com/mateothegreat/go-multilog/multilog"
)

// ReportFn is called to report export progress.
type ReportFn func(obj map[string]interface{})

// Reporter manages when to call the progress callback to reduce noise.
type Reporter struct {
	Fn        ReportFn
	Interval  time.Duration
	BatchSize int

	last  time.Time
	count int
}

// Report calls the progress callback and updates tracking state.
func (pd *Reporter) Report(obj map[string]interface{}) {
	if pd == nil {
		return
	}

	pd.count++

	if time.Since(pd.last) >= pd.Interval {
		pd.Fn(obj)
		pd.last = time.Now()
	}
}

// NewReporter creates a new progress debouncer.
func NewReporter(fn ReportFn, interval time.Duration, batchSize int) *Reporter {
	// Use default debug function if none provided
	if fn == nil {
		fn = func(obj map[string]interface{}) {
			multilog.Debug("plugin.reporter", "Report", obj)
		}
	}

	return &Reporter{
		Fn:        fn,
		Interval:  interval,
		BatchSize: batchSize,
		last:      time.Now(),
		count:     0,
	}
}
