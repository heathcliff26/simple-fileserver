package fileserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/heathcliff26/simple-fileserver/pkg/filesystem"
	"github.com/heathcliff26/simple-fileserver/pkg/middleware"
)

type SSLConfig struct {
	Enabled          bool
	Certificate, Key string
}

type Fileserver struct {
	SSL    SSLConfig
	server *http.Server
}

func NewFileserver(webroot string, index bool) *Fileserver {
	fs := filesystem.CreateFilesystem(webroot, index)
	server := http.FileServer(fs)
	return &Fileserver{
		server: &http.Server{
			Handler:     middleware.Logging(server),
			ReadTimeout: 10 * time.Second,
		},
	}
}

func (s *Fileserver) UseSSL(cert, key string) {
	s.SSL = SSLConfig{
		Enabled:     true,
		Certificate: cert,
		Key:         key,
	}
}

func (s *Fileserver) ListenAndServe(addr string) error {
	s.server.Addr = addr
	slog.Info("Starting server", slog.String("addr", addr))

	var err error
	if s.SSL.Enabled {
		err = s.server.ListenAndServeTLS(s.SSL.Certificate, s.SSL.Key)

	} else {
		err = s.server.ListenAndServe()
	}

	// This just means the server was closed after running
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Calls http.Server.Shutdown(), see it's description for more details
func (s *Fileserver) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
