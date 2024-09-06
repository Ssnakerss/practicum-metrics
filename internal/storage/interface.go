package storage

import (
	"fmt"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

type DataStorage interface {
	New(p ...string) error

	Read(m *metric.Metric) error
	Write(m *metric.Metric) error

	ReadAll(*([]metric.Metric)) (int, error)
	WriteAll(*([]metric.Metric)) (int, error)
	Truncate() error

	CheckStorage() error

	Close()
}

// Описываем ошибку для типа storage в которую будем упаковывать специфические ошибки
const (
	ConnectionError uint = 10
	TimeoutError    uint = 13
	OtherError      uint = 99
)

type StorageError struct {
	Err         error
	StorageType string
	Method      string
	ErrCode     uint
}

func (se *StorageError) Error() string {
	return fmt.Sprintf("[%s][%s][CODE:%d] %v", se.StorageType, se.Method, se.ErrCode, se.Err)
}

func NewStorageError(stp string, mt string, ecode uint, err error) error {
	return &StorageError{
		Err:         err,
		StorageType: strings.ToUpper(stp),
		Method:      strings.ToUpper(mt),
		ErrCode:     ecode,
	}
}
