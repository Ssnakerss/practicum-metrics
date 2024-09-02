package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

// Функция запускающая методы с повторения при возникновении ошибок

func Write(s string) error {
	time.Sleep(1 * time.Second)
	fmt.Println(s)
	return storage.NewStorageError("tst", "write", storage.OtherError, fmt.Errorf("some error %w", errors.New("rrrr")))
}

// Модуль для тестирования и отладки
func execWithRetry(f func(string) error) func(string) error {
	return func(s string) error {
		err := errors.New("trying to exec")
		retry := 0
		var stErr *storage.StorageError
		//Читаем
		//При ошибке -  пробуем еще раз с задежкой
		for err != nil {
			logger.Log.Info("call read")
			time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)
			//Вызываем основной метод
			err = f(s + strconv.Itoa(retry))

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
			logger.SLog.Warnf("error reporting, retry in %d seconds", flags.RetryIntervals[retry])
		}
		return err
	}
}

func main() {
	// cоздаем логгер ZAP
	// не получится - проолжать не имеет смысла, fatal
	if err := logger.Initialize("DEBUG"); err != nil {
		log.Fatal("FATAL: cannot initialize LOGGER: ", err)
	}
	defer logger.Log.Sync()

	fmt.Printf("exex result: %v\n", execWithRetry(Write)("KUKU"))
}
