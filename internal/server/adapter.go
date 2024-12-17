package server

import (
	"context"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/dtadapter"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

func InitAdapter(ctx context.Context, c *app.ServerConfig) (*dtadapter.Adapter, error) {

	//Создаем хранилище
	//Хранилище должно соответствовать интерфейсу storage.DataStorage
	var st storage.DataStorage

	da := dtadapter.Adapter{}
	var filest *storage.FileStorage

	//Если задан DSN - используем БД в качестве хранилища
	if c.DatabaseDSN != "default" && c.DatabaseDSN != "" {
		st = &storage.DBStorage{}

		//Ставим таймаут 60 секунд
		if err := st.New(ctx, c.DatabaseDSN, "60"); err != nil {
			return nil, err
		}

		//Очищаем таблицу
		st.Truncate()
		logger.SLog.Info("using db as storage")
	} else {
		//Иначе используем хранение в памяти
		st = &storage.MemStorage{}
		st.New(ctx)
		logger.SLog.Info("using memory as storage")

		//Если задан путь к файлу - добавляем фаловое хранилище
		if c.StoreFile != "default" {
			filest = &storage.FileStorage{}
			if err := filest.New(context.TODO(), c.StoreFile); err != nil {
				logger.SLog.Warnw("file creation failure", "path", c.StoreFile, "err", err)
			} else {
				if c.Restore {
					//Восстанавливаем значения из файла
					err := da.CopyState(filest, st)
					logger.SLog.Infow("restoring data from ", "file", filest.Filename)
					if err != nil {
						logger.SLog.Warnw("data restore", "failed", err)
					}
				}
			}
		}
	}

	da.New(st, app.RetryIntervals)
	if filest != nil {
		//Очищаем второе хранилище перед записью
		filest.Truncate()
		//Добавляем хранилище и включаем синхронизацию
		//0 - пишем в оба сразе, > 0 - по расписанию
		da.StartSync(c.StoreInterval, filest)
		logger.SLog.Infow("using a sync storage", "file", filest.Filename)
	}
	return &da, nil
}
