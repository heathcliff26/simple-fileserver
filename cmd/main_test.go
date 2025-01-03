package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	if os.Getenv("RUN_CRASH_TEST") == "1" {
		t.Setenv("SFILESERVER_WEBROOT", "../pkg/filesystem/testdata")
		t.Setenv("SFILESERVER_CERT", "ssl.crt")
		t.Setenv("SFILESERVER_KEY", "ssl.key")
		t.Setenv("SFILESERVER_LOG", "true")
		main()
		// Should not reach here, ensure exit with 0 if it does
		os.Exit(0)
	}
	execExitTest(t, "TestMain", true)
}
