package main

import (
	"fmt"
	"time"
)

// Функция запускающая методы с повторения при возникновении ошибок
func execWithRetry(f func(string) error, s string) error {
	fmt.Println("hello from  exec")
	err := f(s)
	if err != nil {
		return err
	}
	return nil
}

func Write(s string) error {
	time.Sleep(3 * time.Second)
	fmt.Println("hello from  write:", s)
	return nil
}

func Read(i int) error {
	time.Sleep(2 * time.Second)
	fmt.Println("hello from Read: ", i)
	return nil
}

// Модуль для тестирования и отладки
func main() {
	execWithRetry(Write, "KUKU")
}
