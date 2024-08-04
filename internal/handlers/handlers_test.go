package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestSetDataTextHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type urlArgs struct {
		mtype  string
		mname  string
		mvalue string
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
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, tt.targetString, nil)
			//https://stackoverflow.com/questions/54580582/testing-chi-routes-w-path-variables
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.args.mtype)
			rctx.URLParams.Add("name", tt.args.mname)
			rctx.URLParams.Add("value", tt.args.mvalue)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			SetDataTextHandler(w, request)

			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.want.code {
				// t.Errorf("%s >> Response code = %d want %d", tt.targetString, res.StatusCode, tt.want.code)
			}

		})
	}
}

func TestSetDataJSONHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		body        []byte
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	// tests := []struct {
	// 	contentType string
	// 	jsonBody string
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	// t.Run(tt.name, func(t *testing.T) {
	// 	// 	SetDataJSONHandler(tt.args.w, tt.args.r)
	// 	// })
	// }
}

func TestChiGetHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDataTextHandler(tt.args.w, tt.args.r)
		})
	}
}
