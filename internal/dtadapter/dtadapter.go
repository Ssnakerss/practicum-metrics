package dtadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

type Adapter struct {
	Ds          storage.DataStorage
	SyncStorage storage.DataStorage
	syncMode    bool
}

func (da *Adapter) New(ds storage.DataStorage) {
	da.Ds = ds
	da.syncMode = false
	da.SyncStorage = nil
}

// Функция-обертка для  повторного вызова при возникновении ошибок
func execRWWtihRetry(f func(*metric.Metric) error) func(*metric.Metric) error {
	return func(m *metric.Metric) error {
		err := errors.New("trying to exec")
		retry := 0
		var stErr *storage.StorageError
		//При ошибке подключения  -  пробуем еще раз с задежкой
		for err != nil {
			logger.Log.Info("call read")
			time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)
			//Вызываем основной метод
			err = f(m)
			//Выходим если нет ошибки или закончился лимит попыток или ошибка не приводится к типу StorageError
			if err == nil ||
				retry == len(flags.RetryIntervals)-1 ||
				!errors.As(err, &stErr) {
				break
			}
			//Если ошибка не связана с подключение к хранилищу -  тоже выходим
			if stErr.ErrCode != storage.ConnectionError {
				break
			}

			retry++
			fmt.Printf("%v\n", err)
			logger.SLog.Warnf("%v, retry in %d seconds", err, flags.RetryIntervals[retry])
		}
		return err
	}
}

// Функция-обертка для  повторного вызова при возникновении ошибок
func execRWAllWtihRetry(f func(*[]metric.Metric) (int, error)) func(*[]metric.Metric) (int, error) {
	return func(m *[]metric.Metric) (int, error) {
		err := errors.New("trying to exec")
		retry := 0
		cnt := 0
		var stErr *storage.StorageError
		//При ошибке подключения  -  пробуем еще раз с задежкой
		for err != nil {
			logger.Log.Info("call read")
			time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)
			//Вызываем основной метод
			cnt, err = f(m)
			//Выходим если нет ошибки или закончился лимит попыток или ошибка не приводится к типу StorageError
			if err == nil ||
				retry == len(flags.RetryIntervals)-1 ||
				!errors.As(err, &stErr) {
				break
			}
			//Если ошибка не связана с подключение к хранилищу -  тоже выходим
			if stErr.ErrCode != storage.ConnectionError {
				break
			}
			retry++
			fmt.Printf("%v\n", err)
			logger.SLog.Warnf("%v, retry in %d seconds", err, flags.RetryIntervals[retry])
		}
		return cnt, err
	}
}

// Пишем в хранилище
// Если интервал синхронизации == 0 - пишем сразу и во второе
func (da *Adapter) Write(m *metric.Metric) error {
	//Вызов записи в базу через ф-ю с повторами при ошибках
	// err := execRWWtihRetry(da.Ds.Write)(m)
	err := da.Ds.Write(m)
	if err != nil {
		return err
	}
	if da.syncMode {
		if err := da.SyncStorage.Write(m); err != nil {
			return err
		}
	}
	return nil
}

func (da *Adapter) WriteAll(mm *[]metric.Metric) error {
	//Вызов записи в базу через ф-ю с повторами при ошибках
	// if _, err := execRWAllWtihRetry(da.Ds.WriteAll)(mm); err != nil {
	if _, err := da.Ds.WriteAll(mm); err != nil {
		errr := fmt.Errorf("data adapter err: %w", err)
		logger.SLog.Error(errr)
		return errr
	}
	//Пишем во второе хранилище если включена синхронная запись
	if da.syncMode {
		if _, err := da.SyncStorage.WriteAll(mm); err != nil {
			return err
		}
	}
	return nil
}

func (da *Adapter) Read(m *metric.Metric) error {
	// return execRWWtihRetry(da.Ds.Read)(m)
	return da.Ds.Read(m)
}

func (da *Adapter) ReadAll(mm *[]metric.Metric) error {
	// _, err := execRWAllWtihRetry(da.Ds.ReadAll)(mm)
	_, err := da.Ds.ReadAll(mm)
	return err
}

// Синхронизация записи
// Если интервал == 0 - синхронная запись во второе хранилище через метод da.Write
// Если интревал > 0 - запускаем горутину с копированием состояния
func (da *Adapter) Sync(interval uint, dst storage.DataStorage) {
	da.SyncStorage = dst
	da.syncMode = (interval == 0)
	if da.syncMode {
		return
	}
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			<-ticker.C
			//Надо почистить перед записью!!!
			da.SyncStorage.Truncate()
			da.CopyState(da.Ds, da.SyncStorage)
		}
	}()
}

// Копирование состояния хранилища
func (da *Adapter) CopyState(src storage.DataStorage, dst storage.DataStorage) error {
	mm := make([]metric.Metric, 0)
	readcnt, err := src.ReadAll(&mm)
	if err != nil {
		return err
	}
	writecnt, err := dst.WriteAll(&mm)
	if err != nil {
		return err
	}
	if readcnt != writecnt {
		return fmt.Errorf("read count %d not equal to write count %d", readcnt, writecnt)
	}

	return nil
}

// Read metric and convert to data interface type
func (da *Adapter) readMetricAndMarshal(m *metric.Metric) ([]byte, error) {
	err := da.Read(m)
	if err != nil {
		return nil, fmt.Errorf("fail to read metric: %w", err)
	}
	mi := metric.ConvertMetricS2I(m)
	mj, err := json.Marshal(mi)
	if err != nil {
		return nil, fmt.Errorf("fail to convert saved metric: %w", err)
	}
	return mj, nil
}

// Cheking if request correct and extract metric from Body
func (da *Adapter) checkRequestAndGetMetric(r *http.Request) (*metric.Metric, error) {
	ct := r.Header.Get("content-type")
	if !strings.Contains(ct, "application/json") {
		return nil, fmt.Errorf("incorrect content type: %v", ct)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read request body: %w", err)
	}
	//Decompression -> TODO: Change to MiddleWare
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			return nil, fmt.Errorf("fail to un-gzip body %w", err)
		}
	}
	var mi metric.MetricJSON
	err = json.Unmarshal(body, &mi)
	if err != nil {
		return nil, fmt.Errorf("fail to convert json: %w", err)
	}
	return metric.ConvertMetricI2S(&mi), nil
}
