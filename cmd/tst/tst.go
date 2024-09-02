package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
)

// Функция запускающая методы с повторения при возникновении ошибок

func Write(s string) error {
	time.Sleep(1 * time.Second)
	fmt.Println(s)
	return fmt.Errorf("read error: %w", errors.New("rrrrr"))
}

// Модуль для тестирования и отладки
func execWithRetry(f func(string) error) func(string) error {
	return func(s string) error {
		err := errors.New("trying to exec")
		retry := 0
		//Читаем
		//При ошибке -  пробуем еще раз с задежкой
		for err != nil {
			logger.Log.Info("call read")
			time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)
			err = f(s + strconv.Itoa(retry))
			if err == nil || retry == len(flags.RetryIntervals)-1 {
				break
			}
			retry++
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

	execWithRetry(Write)("KUKU")
}
