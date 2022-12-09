package main

import (
	"fmt"
	"github.com/verveWireless/vmetrics/vmetrics"
	"time"
)

var (
	personMetric = vmetrics.NewMessage("yield", "bo", "persons")
)

func init() {
	vmetrics.Register(personMetric)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	SomeOperation()

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
