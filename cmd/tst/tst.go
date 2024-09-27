//package for testing code samples

package main

import (
	"context"
	"fmt"
	"time"

	"math/rand"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("рабочий процесс", id, "запущен задача", j)
		time.Sleep(time.Second * 1)
		fmt.Println("рабочий процесс", id, "завершен задача", j)
		results <- j
	}
}

func dummyMetrict(ctx context.Context, readyToReceive chan struct{}) chan []float64 {
	result := make(chan []float64)
	go func() {
		defer close(result)
		for {
			//Задержка для имитации работы процесса
			time.Sleep(time.Second * time.Duration(rand.Intn(3)))

			select {
			case <-ctx.Done():
				return
				//Ждем подтверждения готовности к приему данных
			case <-readyToReceive:
				c, err := cpu.Percent(time.Millisecond*1000, true)
				if err != nil {
					return
				}
				v, err := mem.VirtualMemory()
				if err != nil {
					return
				}
				c = append(c, float64(v.Total), float64(v.Free))
				result <- c
			}
		}
	}()

	return result
}
func gen(ctx context.Context, readyToReceive chan struct{}, cnt int) []chan []float64 {
	res := make([]chan []float64, cnt)
	for i := 0; i < cnt; i++ {
		res[i] = dummyMetrict(ctx, readyToReceive)
	}
	return res
}

func main() {
	const numjobs = 5

	readyToReceive := make(chan struct{})
	defer close(readyToReceive)

	res := gen(context.Background(), readyToReceive, numjobs)

	for i := 0; i < 10; i++ {
		for j := 0; j < numjobs; j++ {
			readyToReceive <- struct{}{}
		}
		for _, c := range res {
			fmt.Println("iter ", i, ">>", <-c)
		}
	}

	// const numjobs = 5
	// jobs := make(chan int, numjobs)
	// results := make(chan int, numjobs)

	// for w := 1; w <= 3; w++ {
	// 	go worker(w, jobs, results)
	// }

	// for j := 1; j <= numjobs; j++ {
	// 	jobs <- j
	// }
	// close(jobs)

	// for a := 1; a <= numjobs; a++ {
	// 	fmt.Println(<-results)
	// }

	// // var wg sync.WaitGroup
	// // var sema = NewSemaphore(2)
	// // for i := 0; i < 10; i++ {
	// // 	wg.Add(1)
	// // 	go func(i int) {
	// // 		defer wg.Done()
	// // 		sema.Acquire()
	// // 		defer sema.Release()
	// // 		fmt.Println(i)
	// // 		time.Sleep(time.Second * 1)
	// // 	}(i)
	// // }
	// // wg.Wait()

}
