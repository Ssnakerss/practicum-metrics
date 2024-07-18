package report

import (
	"log"

	"github.com/go-resty/resty/v2"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func ReportMetrics(mm map[string]metric.Metric, serverAddr string) error {
	for _, m := range mm {
		err := SendMetric(m, serverAddr)
		if err != nil {
			log.Printf("error happened while sending %v: %s \r\n", m, err)
			return err
		}
	}
	return nil
}

func SendMetric(m metric.Metric, serverAddr string) error {
	url := "http://" + serverAddr + "/update/" + m.Type + "/" + m.Name + "/" + m.Value()
	client := resty.New()
	_, err := client.R().Post(url)

	if err != nil {
		return err
	}
	return nil
}
