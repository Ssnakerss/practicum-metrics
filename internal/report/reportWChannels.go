package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/hash"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/go-resty/resty/v2"
)

func SendMetrics(ctx context.Context, mm []metric.Metric) {
	sendChannel := make(chan []metric.Metric) //send channel for metrics
	numWorkers := flags.Cfg.RateLimit
	if numWorkers <= 0 {
		numWorkers = 1
	}
	if numWorkers > len(mm) {
		numWorkers = len(mm)
	}
	//Размер батча - количество метрик для отправки в одном воркере
	batchSize := len(mm) / numWorkers

	//Запускаем воркеров, которые будут отправлять метрики
	createPool(ctx, sendChannel, numWorkers)

	//Отправляем метрики в канал порциями
	for i := 0; i < len(mm); i = i + batchSize {
		end := i + batchSize
		if end > len(mm) {
			end = len(mm)
		}
		sendChannel <- mm[i:end]
	}
	close(sendChannel)
}

func createPool(ctx context.Context, sendChannel <-chan []metric.Metric, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go sendWorker(ctx, sendChannel, i)
	}
}

func sendWorker(ctx context.Context, dataChannel <-chan []metric.Metric, workerNum int) {
	for metrics := range dataChannel {
		logger.SLog.Infof("worker %d start sending %d metrics", workerNum, len(metrics))
		err := SendWithRetry(ctx, metrics, flags.Cfg.EndPointAddress)
		if err != nil {
			logger.SLog.Errorf("worker %d failed to send metrics: %v", workerNum, err)
		} else {
			logger.SLog.Infof("worker %d complete sent %d metrics", workerNum, len(metrics))
		}
	}
}

func SendWithRetry(ctx context.Context, mm []metric.Metric, endpoint string) error {
	//Отправляем метрики
	//При ошибке -  пробуем еще раз с задержкой
	err := errors.New("trying to send")
	retry := 0
	for err != nil {
		//Делаем задержку в случае неудачной попытки отправки метрик
		time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)
		err = ReportMetrics(mm, endpoint)
		//Если удалось отправить
		if err == nil ||
			//или закончилось количество  попыток
			retry == len(flags.RetryIntervals)-1 ||
			//или отмена контекста
			ctx.Err() != nil {
			//то выходим из цикла
			break
		}
		retry++
		logger.SLog.Warnf("send error, retry in %d seconds", flags.RetryIntervals[retry])
	}
	return err
}

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
		return httpSend(body, flags.Cfg.Key, url)
	}
	return nil
}

func httpSend(body []byte, hashKey string, url string) error {
	// посчитаем подпись если задан ключ
	hash, err := hash.MakeSHA256(body, hashKey)
	if err != nil {
		return err
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
