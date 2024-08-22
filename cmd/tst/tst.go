package main

import (
	"fmt"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

func main() {
	f := storage.FileStorage{}
	if err := f.New(`D:\Temp\superfile.txt`); err != nil {
		fmt.Println(err)
		return
	}

	f.Truncate()

	m := metric.Metric{
		Name:    "GetSet30",
		Type:    "counter",
		Counter: 1049634588,
		// Type: "gauge",
		// Gauge: 4.4,
	}

	// fmt.Printf("m original: %v\n", m)

	f.Write(&m)

	f.Read(&m)
	fmt.Printf("m from file: %v\n", m)

	m.Counter = 829353088

	f.Write(&m)

	f.Read(&m)
	fmt.Printf("m from file: %v\n", m)

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
