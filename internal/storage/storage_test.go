package storage

import (
	"testing"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func TestStorage_Insert(t *testing.T) {
	type args struct {
		m metric.Metric
	}
	tests := []struct {
		name    string
		st      *Storage
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.st.Insert(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	type args struct {
		m metric.Metric
	}
	tests := []struct {
		name    string
		st      *Storage
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.st.Update(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
