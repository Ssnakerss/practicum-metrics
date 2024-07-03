package main

import (
	"agent/metric"
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
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

func reportMetrics(mm []metric.Metric) {
	for _, m := range mm {
		err := SendMetric(m)
		if err != nil {
			fmt.Printf("error happened while sending %v: %s \r\n", m, err)
		}
	}
}

func SendMetric(m metric.Metric) error {
	//fmt.Printf("%v\n\r", m)
	url := serverAddr + "update/" + m.MType + "/" + m.Name + "/" + strconv.FormatFloat(m.Value, 'f', -1, 64)
	resp, err := http.Post(url, contentType, bytes.NewReader([]byte(``)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	//fmt.Printf("Status Code: %d\r\n", response.StatusCode)
	resp.Body.Close()
	return nil
}

// func ExtendedPrint(v interface{}) {
// 	val := reflect.ValueOf(v)
// 	// //  проверяем, а не передали ли нам указатель на структуру
// 	// switch val.Kind() {
// 	// case reflect.Ptr:
// 	// 	if val.Elem().Kind() != reflect.Struct {
// 	// 		fmt.Printf("Pointer to %v : %v", val.Elem().Type(), val.Elem())
// 	// 		return
// 	// 	}
// 	// 	// если всё-таки это указатель на структуру, дальше будем работать с самой структурой
// 	// 	val = val.Elem()
// 	// case reflect.Struct: // работаем со структурой
// 	// default:
// 	// 	fmt.Printf("%v : %v", val.Type(), val)
// 	// 	return
// 	// }
// 	fmt.Printf("Struct of type %v and number of fields %d:\n", val.Type(), val.NumField())
// 	for fieldIndex := 0; fieldIndex < val.NumField(); fieldIndex++ {
// 		field := val.Field(fieldIndex) // field — тоже Value
// 		fmt.Printf("\tField %v: %v - val :%v\n", val.Type().Field(fieldIndex).Name, field.Type(), field)
// 		// имя поля мы получаем не из значения поля, а из его типа.
// 	}
// }

func pollMemStatsMetrics(metricsToGather []string, result []metric.Metric) (bool, error) {
	var memoryStat runtime.MemStats
	runtime.ReadMemStats(&memoryStat)
	val := reflect.ValueOf(memoryStat)
	idx := 0
	for _, k := range metricsToGather {

		field := val.FieldByName(k)
		//For better code readability
		var name, value, vtype string
		nm, ok := val.Type().FieldByName(k)
		if ok {
			name = nm.Name
		} else {
			panic("?????")
		}
		value = fmt.Sprintf("%v", field)
		vtype = field.Type().String()
		//----------------------------
		result[idx].Set(name, value, vtype, "gauge")
		idx++
	}
	return true, nil
}

func main() {

	var gatheredMetrics [29]metric.Metric
	//Initialize metrics array for use
	for idx, _ := range gatheredMetrics {
		var m metric.Metric
		m.Set("testgauge", "0", "uint64", "gauge")
		gatheredMetrics[idx] = m
	}

	var cnt uint64 = 0
	rp := 0
	for {

		if rp == reportInterval {
			//It's time to report metrics
			fmt.Print("Reporting metrics ... \r")
			reportMetrics(gatheredMetrics[:])
			rp = 0
		}
		time.Sleep(pollInterval * time.Second)
		rp += pollInterval

		fmt.Printf("%d:Gathering metrics ... \r", cnt)
		_, err := pollMemStatsMetrics(metric.MemStatsMetrics[:], gatheredMetrics[:])
		if err != nil {
			panic(err)
		}
		cnt++
		gatheredMetrics[27].Set("PollCount", strconv.FormatUint(cnt, 10), "uint64", "counter")
		gatheredMetrics[28].Set("RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64), "float64", "gauge")
	}
}
