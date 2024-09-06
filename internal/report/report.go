package report

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func ReportMetrics(mm []metric.Metric, serverAddr string) error {
	//Проверяем есть ли данные для отравки
	if len(mm) > 0 {
		// Для отправки метрик в формате JSON батчами
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
		return httpSend(body, url)
	}
	return nil
}

func httpSend(body []byte, url string) error {
	// посчитаем подпись если задан ключ
	hash := ``
	if flags.Cfg.Key != `` {
		h := hmac.New(sha256.New, []byte(flags.Cfg.Key))
		_, err := h.Write(body)
		if err != nil {
			return err
		}
		hash = hex.EncodeToString(h.Sum(nil))
	}
	//Сжимаем боди гзипом
	bgzip, err := compression.Compress(body)
	if err != nil {
		return err
	}

	client := resty.New()
	_, err = client.R().
		SetHeader("Content-type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("HashSHA256", hash).
		// SetHeader("Accept-Encoding", "gzip").
		SetBody(bgzip).
		Post(url)

	return err
}
