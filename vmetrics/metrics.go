package vmetrics

type Metric interface {
	Record(inf interface{})
	Consume() []string
	Clear()
}

type Opts struct {
	ConstLabels Labels
}
