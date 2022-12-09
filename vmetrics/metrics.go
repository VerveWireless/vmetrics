package vmetrics

type Metric interface {
	Record(inf interface{})
	Consume() []string
	Clear()
	GetName() string
}