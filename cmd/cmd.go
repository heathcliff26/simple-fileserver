package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/heathcliff26/simple-fileserver/pkg/version"
)

const defaultPort = 8080

func init() {
	flag.StringVar(&webroot, "webroot", "", "SFILESERVER_WEBROOT: Required, root directory to serve files from")
	flag.IntVar(&port, "port", defaultPort, "SFILESERVER_PORT: Specify port for the fileserver to listen on")
	flag.StringVar(&sslCert, "cert", "", "SFILESERVER_CERT: SSL certificate to use, needs key as well. Default is no ssl.")
	flag.StringVar(&sslKey, "key", "", "SFILESERVER_KEY: SSL private key to use, needs cert as well. Default is no ssl.")
	flag.BoolVar(&withoutIndex, "no-index", false, "SFILESERVER_NO_INDEX: Do not serve an index for directories, return index.html or 404 instead")
	flag.BoolVar(&enableLogging, "log", false, "SFILESERVER_LOG: Enable logging requests")
	flag.BoolVar(&showVersion, "version", false, "Show the version information and exit")
}

func envBool(target *bool, name string) {
	if val, ok := os.LookupEnv(name); ok {
		*target = strings.ToLower(val) == "true" || val == "1"
	}
}

func envString(target *string, name string) {
	if val, ok := os.LookupEnv(name); ok {
		*target = val
	}
}

// Parse Options not provided by the CLI Arguments from ENV
func parseEnv() {
	if webroot == "" {
		envString(&webroot, "SFILESERVER_WEBROOT")
	}

	if port == defaultPort {
		if val, ok := os.LookupEnv("SFILESERVER_PORT"); ok {
			var err error
			port, err = strconv.Atoi(val)
			if err != nil {
				slog.Error("Could not parse SFILESERVER_PORT", "err", err)
				os.Exit(1)
			}
		}
	}

	if sslCert == "" {
		envString(&sslCert, "SFILESERVER_CERT")
	}

	if sslKey == "" {
		envString(&sslKey, "SFILESERVER_KEY")
	}

	if !withoutIndex {
		envBool(&withoutIndex, "SFILESERVER_NO_INDEX")
	}

	if !enableLogging {
		envBool(&enableLogging, "SFILESERVER_LOG")
	}
}

// Parse CLI Arguments and check the input.
func parseFlags() {
	flag.Parse()
	if showVersion {
		fmt.Print(version.Version())
		os.Exit(0)
	}
	parseEnv()

	if webroot == "" {
		slog.Error("No Webroot: Either -webroot or SFILESERVER_WEBROOT need to be set")
		os.Exit(1)
	}
	if (sslCert != "" && sslKey == "") || (sslCert == "" && sslKey != "") {
		slog.Error("When using ssl need both -cert and -key to be set")
		os.Exit(1)
	}

	slog.Info("Parsed flags",
		slog.String("webroot", webroot),
		slog.Int("port", port),
		slog.String("sslCert", sslCert),
		slog.String("sslKey", sslKey),
		slog.Bool("no-index", withoutIndex),
		slog.Bool("verbose-output", enableLogging),
	)
}
