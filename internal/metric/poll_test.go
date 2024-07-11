package metric

import "testing"

func TestPollMemStatsMetrics(t *testing.T) {
	type args struct {
		metricsToGather []string
		result          map[string]Metric
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
			if err := PollMemStatsMetrics(tt.args.metricsToGather, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("PollMemStatsMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
