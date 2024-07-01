package main

import (
	"bytes"
	"fmt"
	"lib/metric"
	"net/http"
	"strconv"
)

func sendMetric(m metric.Metric) error {
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
