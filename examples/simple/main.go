package main

import (
	"fmt"
	"github.com/VerveWireless/vmetrics/vmetrics"
	"time"
)

var (
	personMetric = vmetrics.NewMessage()
	counterVec   = vmetrics.NewCounterVec([]string{"label1", "label2"})
)

func init() {
	vmetrics.SetupDefaultRegistry([]string{"localhost:9092"}, "phili", nil)
	vmetrics.Register(personMetric)
	vmetrics.Register(counterVec)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	SomeOperation()
	country := "Germany"
	city := "Berlin"
	street := "Karl-Liebknecht"
	counterVec.WithLabelValues(country, city, street).Inc()
	time.Sleep(time.Second * 5)
}

func SomeOperation() {
	fmt.Println("Doing the operation")
	for i := 0; i < 3000; i++ {
		p := Person{Name: "test", Age: i}
		personMetric.Record(p)
		time.Sleep(time.Nanosecond)
	}
}
