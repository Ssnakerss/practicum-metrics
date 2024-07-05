package report

import (
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

const (
	serverAddr  = "http://localhost:8080/"
	contentType = "text/plain"
)

func ReportMetrics(mm []metric.Metric) error {
	for _, m := range mm {
		err := SendMetric(m)
		if err != nil {
			fmt.Printf("error happened while sending %v: %s \r\n", m, err)
			return err
		}
	}
	return nil
}

func SendMetric(m metric.Metric) error {
	// //fmt.Printf("%v\n\r", m)
	// url := serverAddr + "update/" + m.MType + "/" + m.Name + "/" + strconv.FormatFloat(m.Value, 'f', -1, 64)
	// resp, err := http.Post(url, contentType, bytes.NewReader([]byte(``)))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// resp.Body.Close()
	// return nil

	url := serverAddr + "update/" + m.MType + "/" + m.Name + "/" + strconv.FormatFloat(m.Value[0], 'f', -1, 64)
	client := resty.New()
	_, err := client.R().Post(url)

	if err != nil {
		return err
	}

	// if resp.StatusCode() != http.StatusOK {
	// 	//server response bad code
	// }

	return nil

}
