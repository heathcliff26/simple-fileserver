package fileserver

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUseSSL(t *testing.T) {
	fs := NewFileserver("../filesystem/testdata", true)

	assert := assert.New(t)

	assert.Empty(fs.SSL, "New server should not have ssl config yet")

	fs.UseSSL("test.crt", "test.key")

	assert.True(fs.SSL.Enabled, "SSL should be enabled")
	assert.Equal("test.crt", fs.SSL.Certificate, "SSL certificate should match")
	assert.Equal("test.key", fs.SSL.Key, "SSL key should match")
}

func TestListenSSL(t *testing.T) {
	fs := NewFileserver("../filesystem/testdata", true)
	fs.UseSSL("test.crt", "test.key")

	assert := assert.New(t)

	err := fs.ListenAndServe(":8080")

	assert.Error(err, "Should not succeed without valid certificates")
	assert.Equal(":8080", fs.server.Addr, "Should have set addr")
}

func TestListen(t *testing.T) {
	fs := NewFileserver("../filesystem/testdata", true)

	assert := assert.New(t)

	ch := make(chan error, 1)

	go func() {
		err := fs.ListenAndServe(":8080")
		ch <- err
	}()

	res, err := http.Get("http://localhost:8080/test.html")
	assert.NoError(err, "Should be able to reach the server")
	assert.Equal(http.StatusOK, res.StatusCode, "Should receive ok status code")

	assert.NoError(fs.Shutdown(), "Should shutdown server")

	select {
	case err := <-ch:
		assert.NoError(err, "Server should exit without error")
	case <-time.After(time.Second * 5):
		t.Fatal("Server did not shutdown in time")
	}
}
