package report

import (
	"fmt"

	"github.com/go-resty/resty/v2"

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
	url := serverAddr + "update/" + m.Type + "/" + m.Name + "/" + m.Value()
	client := resty.New()
	_, err := client.R().Post(url)

	if err != nil {
		return err
	}

	// if resp.StatusCode() != http.StatusOK {
	// 	//server response bad code
	// }

	return nil

}
