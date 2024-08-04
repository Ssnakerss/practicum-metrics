package report

import (
	"testing"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func TestReportMetrics(t *testing.T) {
	type args struct {
		mm         map[string]metric.Metric
		serverAddr string
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
			if err := ReportMetrics(tt.args.mm, tt.args.serverAddr); (err != nil) != tt.wantErr {
				t.Errorf("ReportMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendMetric(t *testing.T) {
	type args struct {
		m          metric.Metric
		serverAddr string
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
			if err := SendMetric(tt.args.m, tt.args.serverAddr); (err != nil) != tt.wantErr {
				t.Errorf("SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendMetricJSON(t *testing.T) {
	type args struct {
		m          metric.Metric
		serverAddr string
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
			if err := SendMetricJSON(tt.args.m, tt.args.serverAddr); (err != nil) != tt.wantErr {
				t.Errorf("SendMetricJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
