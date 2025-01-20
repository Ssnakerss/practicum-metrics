package subnetchecker

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsInSubnet(t *testing.T) {
	tests := []struct {
		name          string
		subNet        string
		agentIP       string
		result        bool
		wantCreateErr bool
		wantCheckErr  bool
	}{
		{
			name:          "in trusted",
			subNet:        "192.168.0.0/24",
			agentIP:       "192.168.0.1",
			result:        true,
			wantCreateErr: false,
			wantCheckErr:  false,
		},
		{
			name:          "not in trusted",
			subNet:        "192.168.1.0/24",
			agentIP:       "192.168.0.1",
			result:        false,
			wantCreateErr: false,
			wantCheckErr:  false,
		},
		{
			name:    "subnet format err",
			subNet:  "192.168.01/10",
			agentIP: "192.168.0.235",

			wantCreateErr: true,
		},
		{
			name:    "empty subnet",
			subNet:  "",
			agentIP: "192.168.0.235",

			wantCreateErr: true,
		},
		{
			name:    "ip format err",
			subNet:  "192.168.0.1/16",
			agentIP: "192.1680.1",

			wantCreateErr: false,
			wantCheckErr:  true,
		},

		{
			name:    "empty ip",
			subNet:  "192.168.14.0/16",
			agentIP: "",

			wantCreateErr: false,
			wantCheckErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker, err := NewSubNetChecker(tt.subNet)
			if tt.wantCreateErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			res, err := checker.IsTrusted(tt.agentIP)

			if tt.wantCheckErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.result, res)
		})
	}
}

func TestSubNetChecker_Middleware(t *testing.T) {
	tests := []struct {
		name string

		agentIP  string
		wantCode int
	}{
		{
			name:     "in trusted",
			agentIP:  "192.168.0.10",
			wantCode: 200,
		},
		{
			name:     "not in trusted",
			agentIP:  "192.168.10.1",
			wantCode: 403,
		},
		{
			name:     "ip format err",
			agentIP:  "192.168.101",
			wantCode: 500,
		},
		{
			name:     "empty ip",
			agentIP:  "",
			wantCode: 403,
		},
	}

	logger.Initialize("DEBUG")
	defer logger.Log.Sync()

	s, err := NewSubNetChecker("192.168.0.1/24")
	if err != nil {
		log.Fatal("test fails to start")
	}

	testHandler := func() http.HandlerFunc {
		fn := func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
		}
		return fn
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(s.Middleware)
			r.Post("/test", testHandler())

			request := httptest.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()
			request.Header.Add("X-Real-IP", tt.agentIP)
			r.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantCode, res.StatusCode)

		})
	}
}

func TestGetLocalIP(t *testing.T) {
	t.Log(GetLocalIP())
}
