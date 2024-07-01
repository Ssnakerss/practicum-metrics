package main

import (
	"fmt"
	"lib/metric"
	"time"
)

const (
	pollInterval   = 2
	reportInterval = 10
	serverAddr     = "http://localhost:8080/"
	contentType    = "text/plain"
)


var gatheredMetrics []metric.Metric

func reportMetrics(mm *[]metric.Metric) {
	for _, m := range *mm {
		err := sendMetric(m)
		if err != nil {
			fmt.Printf("error happened while sending %v: %s \r\n", m, err)
		}
	}
}

func main() {
	fmt.Println("Agent started ... ")
	cnt := 0
	for {
		fmt.Println("Gathering metrics ... ")
		if cnt == reportInterval {
			fmt.Println("Reporting metrics ... ")
			reportMetrics(&gatheredMetrics)
			cnt = 0
		}
		_, err := pollMetrics(&gatheredMetrics)
		if err != nil {
			panic(err)
		}
		time.Sleep(pollInterval * time.Second)
		cnt += pollInterval
	}
}
