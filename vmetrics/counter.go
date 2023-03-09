package vmetrics

import (
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Counter struct {
	name  string
	value int64
}

func (c *Counter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Get() int64 {
	return c.value
}

func NewCounter(name string) *Counter {
	if strings.Contains(name, "^") {
		panic("Counter name can't contain '^' character")
	}
	return &Counter{name: name, value: 0}
}

func (c *Counter) Aggregated() []string {
	if c.value == 0 {
		return []string{}
	}
	str := "{ \"__time__\": " + strconv.FormatInt(time.Now().UnixNano(), 10) + " , \"name\": \"" + c.name + "\", \"count\": " + strconv.FormatInt(c.value, 10) + "}"
	return []string{str}
}

func (c *Counter) Clear() {
	c.value = 0
}

type CounterVec struct {
	labels   []string
	counters []*Counter
}

func NewCounterVec(labels []string) *CounterVec {
	for _, lbl := range labels {
		if strings.Contains(lbl, "^") {
			panic("Counter Label can't contain '^' character")
		}
	}
	return &CounterVec{
		labels:   labels,
		counters: make([]*Counter, 0),
	}
}

func (cv *CounterVec) GetSize() int {
	return len(cv.counters)
}

func (cv *CounterVec) getCounterWithLabelValues(lvs ...string) *Counter {
	cname := ""
	for _, lbl := range lvs {
		cname += lbl + "^"
	}
	for _, counter := range cv.counters {
		if counter.name == cname {
			return counter
		}
	}
	newCounter := Counter{
		name:  cname,
		value: 1,
	}
	cv.counters = append(cv.counters, &newCounter)
	return &newCounter
}

func (cv *CounterVec) WithLabelValues(lvs ...string) *Counter {
	return cv.getCounterWithLabelValues(lvs...)
}

func (cv *CounterVec) Aggregated() []string {
	var messages []string
	for _, counter := range cv.counters {
		str := "{ \"__time__\": " + strconv.FormatInt(time.Now().UnixNano(), 10) + ","

		labels := strings.Split(counter.name, "^")
		for i, lbl := range cv.labels {
			str += "\"" + lbl + "\": \"" + labels[i] + "\","
		}

		str += "\"count\": " + strconv.FormatInt(counter.value, 10) + "}"
		messages = append(messages, str)
	}
	return messages
}

func (cv *CounterVec) Clear() {
	cv.counters = nil
	cv.counters = make([]*Counter, 0)
}
