package filesystem

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFilesystem(t *testing.T) {
	testMatrix := []struct {
		Name         string
		Path         string
		Index        bool
		ExpectedType string
	}{
		{
			Name:         "withIndex",
			Path:         "foo",
			Index:        true,
			ExpectedType: "http.Dir",
		},
		{
			Name:         "withoutIndex",
			Path:         "foo",
			Index:        false,
			ExpectedType: "filesystem.IndexlessFilesystem",
		},
	}

	for _, tCase := range testMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			fs := CreateFilesystem(tCase.Path, tCase.Index)

			assert.Equal(t, tCase.ExpectedType, reflect.TypeOf(fs).String())
		})
	}
}

func TestIndexlessFilesystem(t *testing.T) {
	fs := IndexlessFilesystem{http.Dir("./testdata")}

	t.Run("DirWithoutIndexFile", func(t *testing.T) {
		assert := assert.New(t)

		f, err := fs.Open("/")
		assert.Nil(f)
		if assert.Error(err) {
			switch osString := strings.ToLower(runtime.GOOS); osString {
			case "windows":
				assert.ErrorContains(err, "The system cannot find the file specified")
			case "linux":
				assert.ErrorContains(err, "open testdata/index.html: no such file or directory")
			default:
				t.Fatalf("Unknown OS %s", osString)
			}
		}
	})
	t.Run("InvalidPath", func(t *testing.T) {
		assert := assert.New(t)

		f, err := fs.Open("not-a-valid-path")
		assert.Nil(f)
		assert.Error(err)
	})

	testMatrix := map[string]string{
		"File":             "/test.html",
		"DirWithIndexFile": "/testdir",
	}

	for name, path := range testMatrix {
		t.Run(name, func(t *testing.T) {
			f, err := fs.Open(path)
			if assert.NoError(t, err) {
				t.Cleanup(func() {
					f.Close()
				})
			}
		})
	}
}

func TestIndexedFilesystem(t *testing.T) {
	fs := http.Dir("./testdata")

	assert := assert.New(t)

	f, err := fs.Open("/test.html")
	assert.Nil(err)
	if assert.NotEmpty(f) {
		f.Close()
	}

	_, err = fs.Open("/nothing.html")
	assert.Error(err)
}
