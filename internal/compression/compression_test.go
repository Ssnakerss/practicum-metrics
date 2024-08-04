package compression

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_gzipWriter_Write(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		w       gzipWriter
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("gzipWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("gzipWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGzipHandle(t *testing.T) {
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GzipHandle(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GzipHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompress(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecompress(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decompress(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decompress() = %v, want %v", got, tt.want)
			}
		})
	}
}
