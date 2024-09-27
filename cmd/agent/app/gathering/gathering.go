package gathering

import (
	"context"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func GatherMetrics(ctx context.Context, pollInterval int, metricType string) <-chan []metric.Metric {
	result := make(chan []metric.Metric)
	go func() {
		pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
		pollCount := 0
		for {
			select {
			case <-ctx.Done():
				close(result)
				return
			case <-pollTicker.C:
				gatheredMetrics := make([]metric.Metric, 0)
				switch metricType {
				case "memstat":
					logger.SLog.Info("Gathering MemStatsMetrics")
					if err := metric.PollMemStatsMetrics(metric.MemStatsMetrics, &gatheredMetrics); err != nil {
						logger.SLog.Errorw("polling metrics", "erorr", err)
					}
				case "extra":
					logger.SLog.Info("Gathering ExtraMetrics")
					pollCount++
					for n, p := range metric.ExtraMetrics {
						var m metric.Metric
						m.Set(n, p.MFunc(pollCount), p.MType)
						gatheredMetrics = append(gatheredMetrics, m)
					}
				case "gops":
					logger.SLog.Info("Gathering Gops")
					if err := metric.PollGopsMetrics(metric.GopsMetrics, &gatheredMetrics); err != nil {
						logger.SLog.Errorw("polling metrics", "erorr", err)
					}
				}
				result <- gatheredMetrics
			}
		}
	}()
	return result
}
