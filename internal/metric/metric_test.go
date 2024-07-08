package metric

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	type args struct {
		mType  string
		mValue string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.args.mType, tt.args.mValue); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_Set(t *testing.T) {
	type args struct {
		mName  string
		mValue string
		mType  string
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
			got, err := tt.m.Set(tt.args.mName, tt.args.mValue, tt.args.mType)
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
