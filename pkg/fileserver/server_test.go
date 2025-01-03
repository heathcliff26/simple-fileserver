package fileserver

import (
	"testing"

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
