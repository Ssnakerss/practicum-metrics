package report

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func ReportMetrics(mm map[string]metric.Metric, serverAddr string) error {
	for _, m := range mm {
		err := SendMetricJSON(m, serverAddr)
		if err != nil {
			log.Printf("error happened while sending %v: %s \r\n", m, err)
			return err
		}
	}
	return nil
}

// Для отправки метрик в текстовом виде
func SendMetric(m metric.Metric, serverAddr string) error {
	url := "http://" + serverAddr + "/update/" + m.Type + "/" + m.Name + "/" + m.Value()
	client := resty.New()
	_, err := client.R().
		SetHeader("content-type", "text/plain").
		Post(url)
	return err
}

// Для отправки метрик в формате JSON
func SendMetricJSON(m metric.Metric, serverAddr string) error {
	url := "http://" + serverAddr + "/update/"

	//Сконвертим метрику в новый формат
	mi := metric.ConvertMetricS2I(&m)
	b, err := json.Marshal(mi)
	logger.SLog.Infow("convert metric", "byte[]", b, "json", string(b))

	if err == nil {
		client := resty.New()
		_, err = client.R().
			SetHeader("Content-type", "application/json").
			SetBody(b).
			Post(url)
	}
	return err
}
