package main

import (
	"fmt"
	"net/http"
	"server/metric"
	"strings"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric

	w.Header().Set("Content-Type", "text/plain")
	//Only POST method allowed
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//Parsing URL to get parameters
	params := strings.Split(r.URL.Path, "/")
	if len(params) != 5 {
		//URL has extra or less parameters
		w.WriteHeader(http.StatusNotFound)
	} else {
		//Checking metric type and name
		if ok, _ := m.Set(params[3], params[4], "float64", params[2]); ok {
			//Processing metrics values
			switch m.MType {
			case "gauge":
				err := Stor.Update(m)
				if err != nil {
					panic("UPDATE ERROR")
				}
			case "counter":
				err := Stor.Insert(m)
				if err != nil {
					panic("INSERT ERROR")
				}
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Storage: \r\n %v", Stor)
	res.Write([]byte(body))
}
