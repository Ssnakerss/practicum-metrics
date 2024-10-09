//package for testing code samples

package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"math/rand"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

type Semaphore struct {
	semaCh chan struct{}
}

func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.semaCh
}

// --------------------------------------------------

func sendWorker(m metric.Metric, done chan struct{}, cnt int) error {
	defer func() {
		<-done
	}()
	//Задержка для имитации работы процесса
	fmt.Printf("worker %d start sending\n\r	", cnt)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	fmt.Printf("by worker %d send metric %v \n\r	", cnt, m)
	atomic.AddInt64(&SentMetrics, 1)
	return nil
}

func anotherWorker(tasks chan metric.Metric) error {
	for m := range tasks {
		fmt.Println("another worker got metric ", m)
	}
	return nil
}

var SentMetrics int64

func main() {
	res := metric.CollectMetrics(1)
	fmt.Println("metrics collected ", len(res))
	// //теперь отправляем полученные метрики
	// //семафор и воркеры
	// const numjobs = 5
	// SentMetrics = 0
	// jobs := make(chan struct{}, numjobs)
	// for _, m := range res {
	// 	jobs <- struct{}{}
	// 	go sendWorker(m, jobs, len(jobs))
	// }
	// for len(jobs) != 0 {
	// 	//ждем пока все задачи выполнятся
	// }
	// close(jobs)
	// fmt.Println("metrics sent ", SentMetrics)
	// //отправим метрики через воркепул
	// //канал для воркеров
	// tasksChannel := make(chan metric.Metric)
	// for i := 0; i < numjobs; i++ {
	// 	go anotherWorker(tasksChannel)
	// }
	// for _, m := range res {
	// 	tasksChannel <- m
	// }
	// close(tasksChannel) //закрыем канал
	//------------------------------------------------
	numWorkers := 1
	batchSize := len(res) / numWorkers
	fmt.Println("batch size ", batchSize)
	for i := 0; i < len(res); i = i + batchSize {
		end := i + batchSize
		if end > len(res) {
			end = len(res)
		}
		fmt.Println("sending batch ", res[i:end])
		// go sendBatch(sample[i:end])
	}
}
