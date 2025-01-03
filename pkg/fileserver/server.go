package fileserver

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/heathcliff26/simple-fileserver/pkg/filesystem"
)

type SSLConfig struct {
	Enabled          bool
	Certificate, Key string
}

type Fileserver struct {
	SSL    SSLConfig
	server http.Handler
}

func NewFileserver(webroot string, index bool) *Fileserver {
	fs := filesystem.CreateFilesystem(webroot, index)
	server := http.FileServer(fs)
	return &Fileserver{
		server: server,
	}
}

type StatusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (rw *StatusResponseWriter) WriteHeader(statusCode int) {
	rw.Status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (s *Fileserver) loggingWrapper(res http.ResponseWriter, req *http.Request) {
	srw := &StatusResponseWriter{ResponseWriter: res}
	s.server.ServeHTTP(srw, req)
	slog.Debug(fmt.Sprintf("Received Request: source=\"%s\", status=%d, path=\"%s\"\n", ReadUserIP(req), srw.Status, req.RequestURI))
}

func (s *Fileserver) Handle(path string) {
	http.HandleFunc("/", s.loggingWrapper)
}

func (s *Fileserver) UseSSL(cert, key string) {
	s.SSL = SSLConfig{
		Enabled:     true,
		Certificate: cert,
		Key:         key,
	}
}

func (s *Fileserver) ListenAndServe(addr string) error {
	slog.Info("Starting server", slog.String("addr", addr))

	var err error
	if s.SSL.Enabled {
		err = http.ListenAndServeTLS(addr, s.SSL.Certificate, s.SSL.Key, nil)

	} else {
		err = http.ListenAndServe(addr, nil)
	}

	// This just means the server was closed after running
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
