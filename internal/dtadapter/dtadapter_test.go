package dtadapter

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestSetGetDataTextHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type urlArgs struct {
		mtype  string
		mname  string
		mvalue string
		method string
	}
	tests := []struct {
		name         string
		targetString string
		args         urlArgs
		want         want
	}{
		{
			name:         "known metric OK test",
			targetString: "/update/{type}/{name}/{value}",
			args: urlArgs{
				mtype:  "gauge",
				mname:  "testM",
				mvalue: "1.1",
				method: "POST",
			},
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "",
			},
		},
		{
			name:         "un-known metric NG test",
			targetString: "/update/{type}/{name}/{value}",
			args: urlArgs{
				mtype:  "unknown",
				mname:  "testM",
				mvalue: "1.1",
				method: "POST",
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "",
			},
		},
		//Checking previously saved metric
		{
			name:         "get known metric",
			targetString: "/value/{type}/{name}",
			args: urlArgs{
				mtype:  "gauge",
				mname:  "testM",
				mvalue: "",
				method: "GET",
			},
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "",
			},
		},
		{
			name:         "get unknown metric",
			targetString: "/value/{type}/{name}",
			args: urlArgs{
				mtype:  "gauge",
				mname:  "testMMM",
				mvalue: "",
				method: "GET",
			},
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
	}

	memst := &storage.MemStorage{}
	if err := memst.New(context.TODO()); err != nil {
		logger.SLog.Fatalw("cannot initialize", "storage", memst)
	}

	da := Adapter{}
	da.New(memst)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			var request *http.Request
			request = httptest.NewRequest(tt.args.method, tt.targetString, nil)
			//https://stackoverflow.com/questions/54580582/testing-chi-routes-w-path-variables
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.args.mtype)
			rctx.URLParams.Add("name", tt.args.mname)
			if tt.args.mvalue != "" {
				rctx.URLParams.Add("value", tt.args.mvalue)
			}

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			switch tt.args.method {
			case "POST":
				da.SetDataTextHandler(w, request)
			case "GET":
				da.GetDataTextHandler(w, request)
			}
			res := w.Result()
			defer res.Body.Close()
			require.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}

func TestSetDataJSONHandler(t *testing.T) {
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

	t.Run("send_json", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		buf.Write([]byte(requestBody))
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.JSONEq(t, requestBody, string(b))
	})

}
