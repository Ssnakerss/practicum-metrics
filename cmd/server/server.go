package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/handlers"
	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

// Middleware для оборачивания хэндлера и логирования событий
// Для обработки запросов
func WithLogging(h http.Handler) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		//
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		//

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
		sugar.Infoln(
			"status", responseData.status,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}

// Теперь займемся ответами
type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter // оигинальный вритер
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

//

func main() {

	// cоздаем логгер ZAP
	// не получится - проолжать не имеет смысла, fatal
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("FATAL: cannot initialize zap")
	}
	defer zapLogger.Sync()
	sugar = *zapLogger.Sugar()

	//If any panic happened during opeartion
	defer func() {
		if err := recover(); err != nil {
			sugar.Fatalf(
				"error heppened while operating -> program will exit",
				"error", err)
		}
	}()

	endPointAddress := ""
	//переменные окружения
	//ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
	if endPointAddress = os.Getenv("ADDRESS"); endPointAddress == "" {
		//Не нашли переменные окружения
		//Параметры командной строки
		//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
		ep := flag.String("a", "localhost:8080", "endpoint address")
		flag.Parse()
		endPointAddress = *ep
		sugar.Infow(
			"use CMD or DEFAULT for config",
			"endPointAddress", endPointAddress,
		)
	} else {
		sugar.Infow(
			"use ENV VAR for config",
			"ADDRESS", endPointAddress,
		)
	}
	//Configuring CHI
	r := chi.NewRouter()
	r.Get("/", WithLogging(http.HandlerFunc(handlers.MainPage)))
	r.Get("/value/{type}/{name}", WithLogging(http.HandlerFunc(handlers.ChiGetHandler)))
	r.Post("/update/{type}/{name}/{value}", WithLogging(http.HandlerFunc(handlers.ChiUpdateHandler)))

	sugar.Infow(
		"starting server at",
		"addr", endPointAddress,
	)

	err = http.ListenAndServe(endPointAddress, r)
	if err != nil {
		sugar.Fatalf(
			"failed to start server -> program will exit",
			"address", endPointAddress,
			"error", err,
		)
	}
}
