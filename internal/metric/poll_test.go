package metric

import "testing"

func TestPollMemStatsMetrics(t *testing.T) {
	type args struct {
		metricsToGather []string
		result          []Metric
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
			got, err := PollMemStatsMetrics(tt.args.metricsToGather, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("PollMemStatsMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PollMemStatsMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
