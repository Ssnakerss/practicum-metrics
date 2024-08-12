package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestSetDataTextHandler(t *testing.T) {
// 	type want struct {
// 		code        int
// 		response    string
// 		contentType string
// 	}
// 	type urlArgs struct {
// 		mtype  string
// 		mname  string
// 		mvalue string
// 	}
// 	tests := []struct {
// 		name         string
// 		targetString string
// 		args         urlArgs
// 		want         want
// 	}{
// 		{
// 			name:         "known metric OK test",
// 			targetString: "/update/{type}/{name}/{value}",
// 			args: urlArgs{
// 				mtype:  "gauge",
// 				mname:  "testM",
// 				mvalue: "1.1",
// 			},
// 			want: want{
// 				code:        http.StatusOK,
// 				response:    "",
// 				contentType: "",
// 			},
// 		},
// 		{
// 			name:         "un-known metric NG test",
// 			targetString: "/update/{type}/{name}/{value}",
// 			args: urlArgs{
// 				mtype:  "unknown",
// 				mname:  "testM",
// 				mvalue: "1.1",
// 			},
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				response:    "",
// 				contentType: "",
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			request := httptest.NewRequest(http.MethodPost, tt.targetString, nil)
// 			//https://stackoverflow.com/questions/54580582/testing-chi-routes-w-path-variables
// 			rctx := chi.NewRouteContext()
// 			rctx.URLParams.Add("type", tt.args.mtype)
// 			rctx.URLParams.Add("name", tt.args.mname)
// 			rctx.URLParams.Add("value", tt.args.mvalue)

// 			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

// 			SetDataTextHandler(w, request)

// 			res := w.Result()
// 			defer res.Body.Close()
// 			if res.StatusCode != tt.want.code {
// 				t.Errorf("%s >> Response code = %d want %d", tt.targetString, res.StatusCode, tt.want.code)
// 			}

// 		})
// 	}
// }

func TestSetDataJSONHandler(t *testing.T) {
	handler := http.HandlerFunc(SetDataJSONHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{
  		"id": "GCSys",
  		"type": "gauge",
  		"value": 2563400
	}`

	// successBody := `{"id":"testSetGet",
	// 				"type":"gauge",
	// 				"value":"4.4"
	// }`

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
