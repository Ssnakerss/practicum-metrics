package dtadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

type Adapter struct {
	Ds          storage.DataStorage
	SyncStorage storage.DataStorage
	syncMode    bool
	ri          []time.Duration
}

func (da *Adapter) New(ds storage.DataStorage, ri []time.Duration) {
	da.Ds = ds
	da.syncMode = false
	da.SyncStorage = nil
	da.ri = ri

}

//TODO: переписать execRWWithRetry & execRWAllWithRetry в одну функцию
//Заменить Read(m *metric.Metric) -> Read(mm *[]metric.Metric),  избавиться от ReadAll/WriteAll

// wrapper functino to retry function call in case of error
// handles Read-Write methods for single metric
func (da *Adapter) execRWWtihRetry(f func(*metric.Metric) error) func(*metric.Metric) error {
	return func(m *metric.Metric) error {
		err := errors.New("trying to exec")
		retry := 0
		var stErr *storage.StorageError
		//При ошибке подключения  -  пробуем еще раз с задежкой
		for err != nil {
			time.Sleep(time.Duration(da.ri[retry]) * time.Second)
			//Вызываем основной метод
			err = f(m)
			//Выходим если нет ошибки или закончился лимит попыток или ошибка не приводится к типу StorageError
			if err == nil ||
				retry == len(da.ri)-1 ||
				!errors.As(err, &stErr) {
				break
			}
			//Если ошибка не связана с подключение к хранилищу -  тоже выходим
			if stErr.ErrCode != storage.ConnectionError {
				break
			}

			retry++
			fmt.Printf("%v\n", err)
			logger.SLog.Warnf("%v, retry in %d seconds", err, da.ri[retry])
		}
		return err
	}
}

// wrapper functino to retry function call in case of error
// handles ReadAll - WriteAll methods for metric slice
func (da *Adapter) execRWAllWtihRetry(f func(*[]metric.Metric) (int, error)) func(*[]metric.Metric) (int, error) {
	return func(m *[]metric.Metric) (int, error) {
		err := errors.New("trying to exec")
		retry := 0
		cnt := 0
		var stErr *storage.StorageError
		//При ошибке подключения  -  пробуем еще раз с задежкой
		for err != nil {
			time.Sleep(time.Duration(da.ri[retry]) * time.Second)
			//Вызываем основной метод
			cnt, err = f(m)
			//Выходим если нет ошибки или закончился лимит попыток или ошибка не приводится к типу StorageError
			if err == nil ||
				retry == len(da.ri)-1 ||
				!errors.As(err, &stErr) {
				break
			}
			//Если ошибка не связана с подключение к хранилищу -  тоже выходим
			if stErr.ErrCode != storage.ConnectionError {
				break
			}
			retry++
			fmt.Printf("%v\n", err)
			logger.SLog.Warnf("%v, retry in %d seconds", err, da.ri[retry])
		}
		return cnt, err
	}
}

// Write single metric to storage
// IF sync interval ==0 - write to sync sto rage
func (da *Adapter) Write(m *metric.Metric) error {
	//Вызов записи в базу через ф-ю с повторами при ошибках
	err := da.execRWWtihRetry(da.Ds.Write)(m)
	// err := da.Ds.Write(m)
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

// Write slice of  metrics to storage
// IF sync interval ==0 - write to sync sto rage
func (da *Adapter) WriteAll(mm *[]metric.Metric) error {
	//Вызов записи в базу через ф-ю с повторами при ошибках
	if _, err := da.execRWAllWtihRetry(da.Ds.WriteAll)(mm); err != nil {
		// if _, err := da.Ds.WriteAll(mm); err != nil {
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

//Read single metric from storage

func (da *Adapter) Read(m *metric.Metric) error {
	return da.execRWWtihRetry(da.Ds.Read)(m)
}

// Read slice of  metrics from storage
func (da *Adapter) ReadAll(mm *[]metric.Metric) error {
	_, err := da.execRWAllWtihRetry(da.Ds.ReadAll)(mm)
	return err
}

// Write sync
// If sync interval == 0 - write to sync storage using da.Write
// If sync interval > 0 - start gorouting to copy state
func (da *Adapter) StartSync(interval uint, dst storage.DataStorage) {
	da.SyncStorage = dst
	da.syncMode = (interval == 0)
	if da.syncMode {
		return
	}
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			<-ticker.C
			da.DoSync()
		}
	}()
}

// Perfom data sync between  storages
func (da *Adapter) DoSync() {
	if da.SyncStorage != nil {
		//Надо почистить перед записью!!!
		da.SyncStorage.Truncate()
		logger.Log.Info("saving state to sync storage")
		da.CopyState(da.Ds, da.SyncStorage)
	}
}

// Copy storage state to another storage
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

	var mi metric.MetricJSON
	err = json.Unmarshal(body, &mi)
	if err != nil {
		return nil, fmt.Errorf("fail to convert json: %w", err)
	}

	return metric.ConvertMetricI2S(&mi), nil
}
