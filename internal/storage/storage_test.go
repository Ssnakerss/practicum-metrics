package storage

import (
	"context"
	"testing"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/stretchr/testify/require"
)

func TestStorage_CounterWriteRead(t *testing.T) {
	f := FileStorage{}
	f.New(context.TODO(), "testfile.txt")
	f.Truncate()
	// f := MemStorage{}
	// f.New()

	t.Run("counter write read test", func(t *testing.T) {
		cm := metric.Metric{
			Name:    "testCounter",
			Type:    "counter",
			Counter: 10,
		}

		//Пишем и проверяем что записалось
		f.Write(&cm)
		rcm := metric.Metric{
			Name: cm.Name,
			Type: cm.Type,
		}

		f.Read(&rcm)
		require.Equal(t, cm, rcm)

		//Пишем еще раз counter и проверяем что значение увеличилось
		f.Write(&cm)
		f.Read(&rcm)

		cm.Counter += cm.Counter

		require.Equal(t, cm, rcm)
	})

}

func TestStorage_GaugeWriteRead(t *testing.T) {
	// f := FileStorage{}
	// f.New("testfile.txt")
	// f.Truncate()
	f := MemStorage{}
	f.New(context.TODO())

	t.Run("gauge write read test", func(t *testing.T) {
		m := metric.Metric{
			Name:  "testGauge",
			Type:  "gauge",
			Gauge: 1.1,
		}
		//Пишем и проверяем что записалось
		f.Write(&m)
		rm := metric.Metric{
			Name: m.Name,
			Type: m.Type,
		}
		f.Read(&rm)
		require.Equal(t, m, rm)

		//Пишем новое значение gauge и проверяем что значение изменилось
		m.Gauge = 55.55
		f.Write(&m)
		f.Read(&rm)
		require.Equal(t, m, rm)
	})
}
