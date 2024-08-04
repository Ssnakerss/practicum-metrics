package main

import (
	"log"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"
)

func main() {

	// cоздаем логгер ZAP
	// не получится - проолжать не имеет смысла, fatal
	if err := logger.Initialize("DEBUG"); err != nil {
		log.Fatal("FATAL: cannot initialize LOGGER: ", err)
	}
	defer logger.Log.Sync()

	m := metric.Metric{
		Name:  "test",
		Type:  "gauge",
		Gauge: 1.1,
	}

	// mj := metric.ConvertMetric(&m)

	// b, err := json.Marshal(mj)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("%v\n\r", b)

	// myString := string(b[:])
	// fmt.Printf("%s\n\r", myString)

	// m = metric.Metric{
	// 	Name:  "test",
	// 	Type:  "counter",
	// 	Gauge: 11000,
	// }

	// fmt.Printf("%v\n\r", m)
	// fmt.Printf("%v\n\r", mj)

	// fmt.Printf("%v -> %v\n\r", mj.Value, *mj.Value)
	// m.Gauge = 2.2
	// fmt.Printf("%v -> %v\n\r", mj.Value, *mj.Value)

	report.SendMetricJSON(m, "localhost:8080")

}
