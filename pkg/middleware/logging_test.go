package middleware

import (
	"bytes"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	t.Cleanup(func() {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	})

	tMatrix := []struct {
		Name     string
		Status   int
		LogLevel slog.Level
		Response string
	}{
		{
			Name:     "DebugLevelOnly",
			LogLevel: slog.LevelInfo,
		},
		{
			Name:     "StatusOK",
			LogLevel: slog.LevelDebug,
			Status:   http.StatusOK,
			Response: "status=200 method=GET path=/test.html",
		},
		{
			Name:     "StatusNotFound",
			LogLevel: slog.LevelDebug,
			Status:   http.StatusNotFound,
			Response: "status=404 method=GET path=/test.html",
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			slog.SetLogLoggerLevel(tCase.LogLevel)

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
				if tCase.Status != 0 {
					res.WriteHeader(tCase.Status)
				}
			})

			handler := Logging(mux)

			req := httptest.NewRequest(http.MethodGet, "/test.html", nil)
			rr := httptest.NewRecorder()

			var buf bytes.Buffer
			log.SetOutput(&buf)
			t.Cleanup(func() {
				log.SetOutput(os.Stderr)
			})

			handler.ServeHTTP(rr, req)
			output := buf.String()

			if tCase.Response == "" {
				assert.Empty(output)
			} else {
				assert.Contains(output, tCase.Response)
			}
		})
	}
}

func TestReadUserIP(t *testing.T) {
	tMatrix := []struct {
		Name, RemoteAddr, RealIP, ForwardedFor, Result string
	}{
		{
			Name:       "RemoteAddr",
			RemoteAddr: "192.168.0.1",
			Result:     "192.168.0.1",
		},
		{
			Name:         "x-forwarded-for",
			ForwardedFor: "192.168.0.2",
			RemoteAddr:   "192.168.0.1",
			Result:       "192.168.0.2",
		},
		{
			Name:         "x-real-ip",
			RealIP:       "192.168.0.3",
			ForwardedFor: "192.168.0.2",
			RemoteAddr:   "192.168.0.1",
			Result:       "192.168.0.3",
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test.html", nil)
			req.RemoteAddr = tCase.RemoteAddr
			req.Header.Set("x-real-ip", tCase.RealIP)
			req.Header.Set("x-forwarded-for", tCase.ForwardedFor)

			assert.Equal(t, tCase.Result, ReadUserIP(req))
		})
	}
}
