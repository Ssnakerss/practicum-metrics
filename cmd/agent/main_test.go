package main

import (
	"lib/metric"
	"testing"
)

func Test_pollMemStatsMetrics(t *testing.T) {
	type args struct {
		metricsToGather []string
		result          []metric.Metric
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pollMemStatsMetrics(tt.args.metricsToGather, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("pollMemStatsMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pollMemStatsMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendMetric(t *testing.T) {
	type args struct {
		m metric.Metric
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
			if err := SendMetric(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
