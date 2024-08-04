package logger

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_loggingResponseWriter_Write(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		r       *loggingResponseWriter
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("loggingResponseWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("loggingResponseWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Initialize(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithLogging(t *testing.T) {
	type args struct {
		h http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLogging(tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLogging() = %v, want %v", got, tt.want)
			}
		})
	}
}
