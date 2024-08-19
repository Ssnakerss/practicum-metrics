package metric

import (
	"testing"
)

func TestPollMemStatsMetrics(t *testing.T) {
	type args struct {
		metricsToGather []string
		err             error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "known metric ",
			args: args{
				metricsToGather: []string{"GCCPUFraction", "HeapAlloc"},
				err:             nil,
			},
			wantErr: false,
		},
		{
			name: "unknown metric ",
			args: args{
				metricsToGather: []string{"GCCPUFraction1", "HeapAlloc1"},
				err:             nil,
			},
			wantErr: true,
		},
		{
			name: "mixed metric ",
			args: args{
				metricsToGather: []string{"GCCPUFraction", "HeapAlloc1"},
				err:             nil,
			},
			wantErr: true,
		},
	}

	gatheredMetrics := make(map[string]Metric)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PollMemStatsMetrics(tt.args.metricsToGather, gatheredMetrics); (err != nil) != tt.wantErr {
				t.Errorf("PollMemStatsMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
