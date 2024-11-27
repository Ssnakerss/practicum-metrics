package dtadapter

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ssnakerss/practicum-metrics/internal/storage"
)

func TestAdapter_MainPage(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		da   *Adapter
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.da.MainPage(tt.args.w, tt.args.r)
		})
	}
}

func ExampleAdapter_SetDataTextHandler() {
	da := Adapter{}
	memst := &storage.MemStorage{}
	memst.New(context.TODO())
	da.New(memst)

	handler := http.HandlerFunc(da.SetDataJSONHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{
  		"id": "GCSys",
  		"type": "gauge",
  		"value": 2563400
	}`

	buf := bytes.NewBuffer(nil)
	buf.Write([]byte(requestBody))
	r := httptest.NewRequest("POST", srv.URL, buf)
	r.RequestURI = ""
	r.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(r)
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	// Output:
	// 200

}
