package storage

import "github.com/Ssnakerss/practicum-metrics/internal/metric"

type DataStorage interface {
	New(p ...string) error

	Read(m *metric.Metric) error
	Write(m *metric.Metric) error

	ReadAll(*([]metric.Metric)) (int, error)
	WriteAll(*([]metric.Metric)) (int, error)
	Truncate() error

	CheckStorage() error
}
