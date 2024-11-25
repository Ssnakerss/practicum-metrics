package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
)

func MakeSHA256(b []byte, key string) (string, error) {
	hash := ``
	if key != `` {
		h := hmac.New(sha256.New, []byte(flags.Cfg.Key))
		_, err := h.Write(b)
		if err != nil {
			return ``, err
		}
		hash = hex.EncodeToString(h.Sum(nil))
	}
	return hash, nil
}

type copyWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (cw copyWriter) Write(b []byte) (int, error) {
	return cw.Writer.Write(b)
}

func HashHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Есди ключ не задан - ничего не делаем
		if flags.Cfg.Key == `` {
			next.ServeHTTP(w, r)
			return
		}

		reqHash := r.Header.Get("HashSHA256")
		//Если задан заголовок - проверяем хэш запроса
		if reqHash != `` {
			//Проверяем хэш запроса
			//Прочитали весь боди запроса - надо потом вернуть обратно
			reqBody, err := io.ReadAll(r.Body)
			//Возвращаем боди обратно в запрос
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			if err != nil {
				//Если ошибка - передаем управление дальше и возвращаемся
				logger.SLog.Errorw("HashHandle read body", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			calcHash, err := MakeSHA256(reqBody, flags.Cfg.Key)
			if err != nil {
				//Если ошибка - передаем управление дальше и возвращаемся
				logger.SLog.Errorw("HashHandle MakeSHA256", "error", err)
				next.ServeHTTP(w, r)
				return
			}
			if string(reqHash) != calcHash {
				//Если хэш не совпадает то выдаем ошибку
				//и прекращаем обработку запроса
				logger.SLog.Warn("Неверный хэш запроса")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		//Обрабатываем ответ
		var body bytes.Buffer
		//Перехватываем w.Write чтобы успеть записать Header
		//Пишем боди в байты чтобы посчитать жэш
		next.ServeHTTP(copyWriter{ResponseWriter: w, Writer: &body}, r)

		hash, err := MakeSHA256(body.Bytes(), flags.Cfg.Key)
		if err != nil {
			logger.SLog.Error(err)
		} else {
			//Пишем хэш в заголовок
			w.Header().Set("HashSHA256", hash)
		}
		//Записываем байты в ответ
		w.Write(body.Bytes())
	})
}
