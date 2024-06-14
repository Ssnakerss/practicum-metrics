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
	metricName  string
	metricType  string //counter || gauge
	metricValue float64
}

func pollMetrics(mm *[]metric) (bool, error) {
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
		switch v := d.(type) {
		case uint64:
			(*mm)[idx].metricValue = float64(d.(uint64))
		case uint32:
			(*mm)[idx].metricValue = float64(d.(uint32))
		case float64:
			(*mm)[idx].metricValue = d.(float64)
		default:
			return false, fmt.Errorf("Unexpected type %v", v)
		}
	}
	(*mm)[27].metricValue += 1
	(*mm)[28].metricValue = rand.Float64()
	return true, nil
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

	resp, err := http.Post(url, contentType, bytes.NewReader([]byte(``)))

	if err != nil {
		fmt.Println(err)
		return err
	}

	//fmt.Printf("Status Code: %d\r\n", response.StatusCode)
	resp.Body.Close()
	return nil
}

func main() {
	fmt.Println("Agent started ... ")
	metricsToGather := []metric{
		{"Alloc", "gauge", 0.0},
		{"BuckHashSys", "gauge", 0.0},
		{"Frees", "gauge", 0.0},
		{"GCCPUFraction", "gauge", 0.0},
		{"GCSys", "gauge", 0.0},
		{"HeapAlloc", "gauge", 0.0},
		{"HeapIdle", "gauge", 0.0},
		{"HeapInuse", "gauge", 0.0},
		{"HeapObjects", "gauge", 0.0},
		{"HeapReleased", "gauge", 0.0},
		{"HeapSys", "gauge", 0.0},
		{"LastGC", "gauge", 0.0},
		{"Lookups", "gauge", 0.0},
		{"MCacheInuse", "gauge", 0.0},
		{"MCacheSys", "gauge", 0.0},
		{"MSpanInuse", "gauge", 0.0},
		{"MSpanSys", "gauge", 0.0},
		{"Mallocs", "gauge", 0.0},
		{"NextGC", "gauge", 0.0},
		{"NumForcedGC", "gauge", 0.0},
		{"NumGC", "gauge", 0.0},
		{"OtherSys", "gauge", 0.0},
		{"PauseTotalNs", "gauge", 0.0},
		{"StackInuse", "gauge", 0.0},
		{"StackSys", "gauge", 0.0},
		{"Sys", "gauge", 0.0},
		{"TotalAlloc", "gauge", 0.0},
		{"PollCount", "counter", 0},   //27
		{"RandomValue", "gauge", 0.0}, //28
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
		_, err := pollMetrics(&metricsToGather)
		if err != nil {
			panic(err)
		}
		time.Sleep(pollInterval * time.Second)
		cnt += pollInterval
	}
}
