package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/tools"
)

type FileStorage struct {
	// DataStorage
	filename string
}

const (
	//TODO: посчитать подходящий размер
	chunckSize int64 = 100
)

func (filest *FileStorage) New(p ...string) error {
	if len(p) < 1 {
		return fmt.Errorf("file name no specified")
	}
	filest.filename = p[0]

	file, err := os.OpenFile(filest.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func (filest *FileStorage) Write(m *metric.Metric) error {
	file, err := os.OpenFile(filest.filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//Находим позицию искомой метрики в файле для обновления записи
	newLine := true
	var pos, i int64
	pos = -1
	i = -1

	//Метрика, которую можем найдем в файле
	mm := metric.Metric{}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		next := true
		for next {
			i++
			l := scanner.Text()
			if err = json.Unmarshal([]byte(l), &mm); err != nil {
				return err
			}

			if (*m).Name+(*m).Type == mm.Name+mm.Type {
				//Нашли нашу метрику
				pos = i
				break
			}
			next = scanner.Scan()
			metric.ClearMetric(&mm)
		}
	}

	//Метрика, которую будем писать
	//Случае новой метрики - это метрика, которую нам передали
	//Если метрика есть в файле -  используем  ее
	rm := m

	//Добавлять ли перевод строки. Если заменяем метрику -  то не надо
	if pos == -1 {
		//Не нашли метрику -  пишем в конец файла
		if _, err := file.Seek(0, 2); err != nil {
			return err
		}
	} else {
		//Нашли метрику - перезаписываем ее
		//Ставим курсор на найденную метрику
		if _, err := file.Seek(pos*(chunckSize+1), 0); err != nil {
			return err
		}
		newLine = false
		//Если метрика Counter - надо увеличить значение на значение из файла

		mm.Counter += (*m).Counter
		mm.Gauge = (*m).Gauge

		rm = &mm
	}

	data, err := json.Marshal(rm)
	if err != nil {
		return err
	}
	//Выравниваем длинну данных
	s := string(data)
	tools.PadRight(&s, " ", chunckSize)
	data = []byte(s)

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if newLine {
		if err := writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

func (filest *FileStorage) WriteAll(mm *([]metric.Metric)) (int, error) {
	cnt := 0
	for _, m := range *mm {
		err := filest.Write(&m)
		if err != nil {
			return cnt, err
		}
		cnt++
	}
	return cnt, nil
}

func (filest *FileStorage) Read(m *metric.Metric) error {
	file, err := os.OpenFile(filest.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return err
	}

	next := true
	mm := metric.Metric{}
	for next {
		l := scanner.Text()
		if err = json.Unmarshal([]byte(l), &mm); err != nil {
			return err
		}
		if m.Name+m.Type == mm.Name+mm.Type {
			m.Counter = mm.Counter
			m.Gauge = mm.Gauge
			return nil
		}
		next = scanner.Scan()
		metric.ClearMetric(&mm)
	}
	return nil
}

func (filest *FileStorage) ReadAll(mm *([]metric.Metric)) (int, error) {
	buf, err := os.ReadFile(filest.filename)
	if err != nil {
		return 0, err
	}
	js := "[" + strings.TrimRight(strings.Replace(string(buf), "\n", ",", -1), ",") + "]"

	err = json.Unmarshal([]byte(js), mm)
	if err != nil {
		return 0, err
	}
	return len(*mm), nil
}

func (filest *FileStorage) Truncate() error {
	if filest.filename != "" {
		f, err := os.OpenFile(filest.filename, os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func (filest *FileStorage) CheckStorage() error {
	return nil
}
