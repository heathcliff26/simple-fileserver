package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/heathcliff26/simple-fileserver/pkg/fileserver"
)

var (
	webroot       string
	port          int
	sslCert       string
	sslKey        string
	withoutIndex  bool
	enableLogging bool
	showVersion   bool
)

func main() {
	parseFlags()

	fs := fileserver.NewFileserver(webroot, !withoutIndex)

	fs.Handle("/")

	if sslCert != "" && sslKey != "" {
		slog.Info("Enabling ssl")
		fs.UseSSL(sslCert, sslKey)
	}

	if enableLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	err := fs.ListenAndServe(":" + strconv.Itoa(port))
	if err != nil {
		slog.Error("Failed to run server", "err", err)
		os.Exit(1)
	}
}
