package fileserver

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingWrapper(t *testing.T) {
	fs := NewFileserver("../filesystem/testdata", true)

	req := httptest.NewRequest(http.MethodGet, "/test.html", nil)
	rr := httptest.NewRecorder()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
	})

	assert := assert.New(t)

	fs.loggingWrapper(rr, req)

	assert.Empty(buf.String())

	fs.Log = true

	fs.loggingWrapper(rr, req)
	output := buf.String()

	assert.Contains(output, "status=200")
	assert.Contains(output, "path=\"/test.html\"")
}

func TestUseSSL(t *testing.T) {
	fs := NewFileserver("../filesystem/testdata", true)

	assert := assert.New(t)

	assert.Empty(fs.SSL, "New server should not have ssl config yet")

	fs.UseSSL("test.crt", "test.key")

	assert.True(fs.SSL.Enabled, "SSL should be enabled")
	assert.Equal("test.crt", fs.SSL.Certificate, "SSL certificate should match")
	assert.Equal("test.key", fs.SSL.Key, "SSL key should match")
}
