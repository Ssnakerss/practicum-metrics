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
		{
			name: "Gauge_OK",
			args: args{
				mType:  "gauge",
				mValue: "10.1",
			},
			want: true,
		},
		{
			name: "Counter_OK",
			args: args{
				mType:  "counter",
				mValue: "10",
			},
			want: true,
		},
		{
			name: "Gauge_NG",
			args: args{
				mType:  "gauge",
				mValue: "1A",
			},
			want: false,
		},
		{
			name: "Counter_NG",
			args: args{
				mType:  "counter",
				mValue: "10.222",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.args.mType, tt.args.mValue); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_Value(t *testing.T) {
	tests := []struct {
		name string
		m    *Metric
		want string
	}{
		{
			name: "Counter_OK",
			m: &Metric{
				Name:    "Counter1",
				Type:    "counter",
				Gauge:   1.1,
				Counter: 2,
			},
			want: "2",
		},
		{
			name: "Gauge_OK",
			m: &Metric{
				Name:    "Gauge1",
				Type:    "gauge",
				Gauge:   1.1,
				Counter: 2,
			},
			want: "1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Value(); got != tt.want {
				t.Errorf("Metric.Value() = %v, want %v", got, tt.want)
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
		wantErr bool
	}{
		{
			name: "Set_NG",
			m: &Metric{
				Name:    "",
				Type:    "",
				Gauge:   0,
				Counter: 0,
			},
			args: args{
				mName:  "nameNG",
				mValue: "1",
				mType:  "unknown",
			},
			wantErr: true,
		},
		{
			name: "Set_OK",
			m: &Metric{
				Name:    "",
				Type:    "",
				Gauge:   0,
				Counter: 0,
			},
			args: args{
				mName:  "CounterOK",
				mValue: "1",
				mType:  "counter",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Set(tt.args.mName, tt.args.mValue, tt.args.mType); (err != nil) != tt.wantErr {
				t.Errorf("Metric.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
