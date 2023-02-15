package vmetrics

import (
	"errors"
	"sync/atomic"
)

type CounterOpts Opts

type Counter interface {
	Metric

	Inc()
	Add(float64)
}

type counter struct {
	value uint64
}

func (c *counter) Add(v float64) {
	if v < 0 {
		panic(errors.New("counter cannot decrease in value"))
	}

	val := uint64(v)
	if float64(val) == v {
		atomic.AddUint64(&c.value, val)
		return
	}
}

func (c *counter) Inc() {
	atomic.AddUint64(&c.value, 1)
}

type CounterVec struct {
	*metricVec
}

func (v *CounterVec) WithLabelValues(lvs ...string) Counter {
	c, err := v.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return c
}

func (v *CounterVec) GetMetricWithLabelValues(lvs ...string) (Counter, error) {
	metric, err := v.metricVec.getMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Counter), err
	}
	return nil, err
}

func NewCounterVec(labelNames []string) *CounterVec {
	return &CounterVec{
		metricVec: newMetricVec(func(lvs ...string) Metric {
			result := &counter{labelPairs: makeLabelPairs(desc, lvs)}
			result.init(result)
			return result
		}),
	}
}
