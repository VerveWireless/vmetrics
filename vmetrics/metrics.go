package vmetrics

type Metric interface {
	Record(inf interface{})
	Consume() []string
	Aggregated() []string
	Clear()
}
