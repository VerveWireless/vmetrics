package vmetrics

type Metric interface {
	Aggregated() []string
	Clear()
}
