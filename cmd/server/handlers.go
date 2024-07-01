package main

import (
	"fmt"
	"lib/metric"
	"net/http"
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
		if !m.IsValid(params[3], params[2]) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			m.Set(params[3], params[4], "float64", params[2])
			//Processing metrics values
			switch m.MType {
			case "gauge":
				Stor.Update(m)
			case "counter":
				Stor.Insert(m)
			}

			w.WriteHeader(http.StatusOK)
		}
	}
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", req.Method)
	body += "Header ===============\r\n"
	for k, v := range req.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Query parameters ===============\r\n"
	for k, v := range req.URL.Query() {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	res.Write([]byte(body))
}
