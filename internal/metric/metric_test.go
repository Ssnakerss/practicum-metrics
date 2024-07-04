package metric

import (
	"testing"
)

func TestMetric_convertValue(t *testing.T) {
	type args struct {
		value string
		vType string
	}
	tests := []struct {
		name    string
		m       *Metric
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.convertValue(tt.args.value, tt.args.vType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.convertValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Metric.convertValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_IsValid(t *testing.T) {
	type args struct {
		name  string
		mType string
	}
	tests := []struct {
		name string
		m    *Metric
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(tt.args.name, tt.args.mType); got != tt.want {
				t.Errorf("Metric.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_Set(t *testing.T) {
	type args struct {
		name  string
		value string
		vType string
		mType string
	}
	tests := []struct {
		name    string
		m       *Metric
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Set(tt.args.name, tt.args.value, tt.args.vType, tt.args.mType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Metric.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}
