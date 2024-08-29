package report

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func ReportMetrics(mm map[string]metric.Metric, serverAddr string) error {
	//Старый формат передачи по одной метрике
	// for _, m := range mm {
	// 	err := SendMetricJSON(m, serverAddr)
	// 	if err != nil {
	// 		log.Printf("error happened while sending %v: %s \r\n", m, err)
	// 		return err
	// 	}
	// }

	//Проверяем есть ли данные для отравки
	if len(mm) > 0 {
		//Отправляем метрики батчами
		if err := SendMetricJSONSlice(mm, serverAddr); err != nil {
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
	if err != nil {
		logger.SLog.Errorw("error marshal", "metric", mi)
		return fmt.Errorf("error marshal metric %v", mi)
	}
	logger.SLog.Infow("convert metric", "byte[]", b, "json", string(b))

	bgzip, err := compression.Compress(b)

	if err == nil {
		client := resty.New()
		_, err = client.R().
			SetHeader("Content-type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			// SetHeader("Accept-Encoding", "gzip").
			SetBody(bgzip).
			Post(url)
	}
	return err
}

// Для отправки метрик в формате JSON батчами
func SendMetricJSONSlice(mm map[string]metric.Metric, serverAddr string) error {
	url := "http://" + serverAddr + "/updates/"
	mcsj := make([]metric.MetricJSON, 0)

	for _, m := range mm {
		//Сконвертим метрику в interface формат
		mi := metric.ConvertMetricS2I(&m)
		mcsj = append(mcsj, *mi)
	}

	body, err := json.Marshal(mcsj)
	if err != nil {
		return fmt.Errorf("error marshal []metricJSON %v", mcsj)
	}
	//Сжимаем боди гзипом
	bgzip, err := compression.Compress(body)

	if err == nil {
		client := resty.New()
		_, err = client.R().
			SetHeader("Content-type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			// SetHeader("Accept-Encoding", "gzip").
			SetBody(bgzip).
			Post(url)
	}
	return err
}
