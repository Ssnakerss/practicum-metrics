package report

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

const (
	serverAddr  = "http://localhost:8080/"
	contentType = "text/plain"
)

func ReportMetrics(mm []metric.Metric) error {
	for _, m := range mm {
		err := SendMetric(m)
		if err != nil {
			fmt.Printf("error happened while sending %v: %s \r\n", m, err)
			return err
		}
	}
	return nil
}

func SendMetric(m metric.Metric) error {
	//fmt.Printf("%v\n\r", m)
	url := serverAddr + "update/" + m.MType + "/" + m.Name + "/" + strconv.FormatFloat(m.Value, 'f', -1, 64)
	resp, err := http.Post(url, contentType, bytes.NewReader([]byte(``)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	resp.Body.Close()
	return nil
}
