package main

import (
	"fmt"
	"strconv"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

func main() {
	m := metric.Metric{
		Name:    "tst0",
		Type:    "counter",
		Gauge:   1,
		Counter: 0,
	}

	f := storage.FileStorage{}
	if err := f.New(`D:\Temp\superfile.txt`); err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 2; i++ {
		m.Counter += 550
		m.Name = "tst" + strconv.Itoa(i)
		fmt.Println(m.Name)
		if err := f.Write(&m); err != nil {
			fmt.Println(err)
		}
	}

	// f.Read(&m)
	// fmt.Printf("m: %v\n", m)

	// var mm []metric.Metric

	// f.ReadAll(&mm)
	// fmt.Printf("m array: %v\n", mm)

	// 	file, err := os.OpenFile(`D:\Temp\superfile.txt`, os.O_RDWR|os.O_CREATE, 0666)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer file.Close()
	// 	file.Seek(0, 2)
	// 	file.Write([]byte("1111111111"))
}
