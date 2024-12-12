package agent

import (
	"testing"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func TestGenMemStats(t *testing.T) {
	for tt := range GenMemStats(metric.MemStatsMetrics) {
		t.Log(tt)
	}
}

func TestGenGopsStats(t *testing.T) {
	for tt := range GenGopsStats() {
		t.Log(tt)
	}
}

func TestGenExtraStats(t *testing.T) {
	for tt := range GenExtraStats(1) {
		t.Log(tt)
	}

}

func TestCollectMetrics(t *testing.T) {
	tt := CollectMetrics(1)
	l := len(tt)
	t.Log(l)
}
