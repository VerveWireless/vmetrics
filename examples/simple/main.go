package main

import (
	"fmt"
	"github.com/VerveWireless/vmetrics/vmetrics"
	"strconv"
	"time"
)

var (
	logMessagesDropped             = vmetrics.NewCounter("LogMessageDropped")
	doubleVerifyRequestStatusCount = vmetrics.NewCounterVec([]string{"status", "ad_format"})
)

func init() {
	vmetrics.SetupDefaultRegistry([]string{"localhost:9092"}, "dv_count_vec", nil)
	vmetrics.MustRegister(doubleVerifyRequestStatusCount, logMessagesDropped)
}

type DVRS struct {
	Status   string `json:"status"`
	AdFormat int    `json:"ad_format"`
}

func main() {
	SomeOperation()
	time.Sleep(time.Second * 20)
}

func SomeOperation() {
	fmt.Println("Doing the operation")
	for i := 0; i < 100; i++ {
		logMessagesDropped.Inc()
		d := DVRS{Status: "blocked_ipv6", AdFormat: 23}
		doubleVerifyRequestStatusCount.WithLabelValues(d.Status, strconv.Itoa(d.AdFormat)).Inc()
		doubleVerifyRequestStatusCount.WithLabelValues(d.Status, "24").Inc()
		doubleVerifyRequestStatusCount.WithLabelValues(d.Status, "25").Inc()
		time.Sleep(time.Nanosecond)
	}
}
