package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const (
	pollInterval   = 2
	reportInterval = 10
	serverAddr     = "http://localhost:8080/"
	contentType    = "text/plain"
)

type metric struct {
	metricName      string
	metricType      string //counter || gauge
	metricValueType string //uint float64
	metricValue     float64
}

func pollMetrics(mm *[]metric) {
	var memoryStat runtime.MemStats
	runtime.ReadMemStats(&memoryStat)

	fields := &memoryStat
	for idx, d := range []interface{}{
		fields.Alloc,
		fields.BuckHashSys,
		fields.Frees,
		fields.GCCPUFraction,
		fields.GCSys,
		fields.HeapAlloc,
		fields.HeapIdle,
		fields.HeapInuse,
		fields.HeapObjects,
		fields.HeapReleased,
		fields.HeapSys,
		fields.LastGC,
		fields.Lookups,
		fields.MCacheInuse,
		fields.MCacheSys,
		fields.MSpanInuse,
		fields.MSpanSys,
		fields.Mallocs,
		fields.NextGC,
		fields.NumForcedGC,
		fields.NumGC,
		fields.OtherSys,
		fields.PauseTotalNs,
		fields.StackInuse,
		fields.StackSys,
		fields.Sys,
		fields.TotalAlloc,
	} {

		switch (*mm)[idx].metricValueType {
		case "uint64":
			(*mm)[idx].metricValue = float64(d.(uint64))
			// fmt.Printf("Index: %d >> value: %#v > %f> Name: %s\r\n", idx, d, float64(d.(uint64)), (*mm)[idx].metricName)
		case "uint32":
			(*mm)[idx].metricValue = float64(d.(uint32))
			// fmt.Printf("Index: %d >> value: %#v > %f> Name: %s\r\n", idx, d, float64(d.(uint32)), (*mm)[idx].metricName)
		case "float64":
			(*mm)[idx].metricValue = d.(float64)
			// fmt.Printf("Index: %d >> value: %#v > %f> Name: %s\r\n", idx, d, d.(float64), (*mm)[idx].metricName)
		}
		//
	}
	(*mm)[27].metricValue += 1
	(*mm)[28].metricValue = rand.Float64()
}
func reportMetrics(mm *[]metric) {
	for _, m := range *mm {
		err := sendMetric(m)
		if err != nil {
			fmt.Printf("Error happened while sending %v: %s \r\n", m, err)
		}
	}
}
func sendMetric(m metric) error {
	url := serverAddr + "update/" + m.metricType + "/" + m.metricName + "/" + strconv.FormatFloat(m.metricValue, 'f', -1, 64)
	data := []byte(``)
	r := bytes.NewReader(data)
	resp, err := http.Post(url, contentType, r)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		//fmt.Printf("Status Code: %d\r\n", response.StatusCode)
		return nil
	}
	defer resp.Body.Close()
}

func main() {
	fmt.Println("Agent started ... ")
	metricsToGather := []metric{
		{"Alloc", "gauge", "uint64", 0.0},
		{"BuckHashSys", "gauge", "uint64", 0.0},
		{"Frees", "gauge", "uint64", 0.0},
		{"GCCPUFraction", "gauge", "float64", 0.0},
		{"GCSys", "gauge", "uint64", 0.0},
		{"HeapAlloc", "gauge", "uint64", 0.0},
		{"HeapIdle", "gauge", "uint64", 0.0},
		{"HeapInuse", "gauge", "uint64", 0.0},
		{"HeapObjects", "gauge", "uint64", 0.0},
		{"HeapReleased", "gauge", "uint64", 0.0},
		{"HeapSys", "gauge", "uint64", 0.0},
		{"LastGC", "gauge", "uint64", 0.0},
		{"Lookups", "gauge", "uint64", 0.0},
		{"MCacheInuse", "gauge", "uint64", 0.0},
		{"MCacheSys", "gauge", "uint64", 0.0},
		{"MSpanInuse", "gauge", "uint64", 0.0},
		{"MSpanSys", "gauge", "uint64", 0.0},
		{"Mallocs", "gauge", "uint64", 0.0},
		{"NextGC", "gauge", "uint64", 0.0},
		{"NumForcedGC", "gauge", "uint32", 0.0},
		{"NumGC", "gauge", "uint32", 0.0},
		{"OtherSys", "gauge", "uint64", 0.0},
		{"PauseTotalNs", "gauge", "uint64", 0.0},
		{"StackInuse", "gauge", "uint64", 0.0},
		{"StackSys", "gauge", "uint64", 0.0},
		{"Sys", "gauge", "uint64", 0.0},
		{"TotalAlloc", "gauge", "uint64", 0.0},
		{"PollCount", "counter", "uint64", 0},   //27
		{"RandomValue", "gauge", "uint64", 0.0}, //28
	}
	//
	// var memoryStat runtime.MemStats
	// runtime.ReadMemStats(&memoryStat)
	// v := reflect.ValueOf(memoryStat)
	// typeOfS := v.Type()
	// for i := 0; i < v.NumField(); i++ {
	// 	val := v.Field(i).Interface().(uint64)
	// 	name := typeOfS.Field(i).Name
	// 	fmt.Println(name, " >>> ", val)
	// }
	// return
	//
	cnt := 0
	for {
		fmt.Println("Gathering metrics ... ")
		if cnt == reportInterval {
			fmt.Println("Reporting metrics ... ")
			reportMetrics(&metricsToGather)
			cnt = 0
		}
		pollMetrics(&metricsToGather)
		time.Sleep(pollInterval * time.Second)
		cnt += pollInterval
	}
}
