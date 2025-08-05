package e2e

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/heathcliff26/simple-fileserver/tests/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifyPages(t *testing.T, port int) {
	assert := assert.New(t)

	baseURL := fmt.Sprintf("http://localhost:%d", port)

	body := fetchAndReadPage(t, baseURL, "/index.html")
	assert.Equal(indexHTML, string(body), "Index page should match expected content")
	body = fetchAndReadPage(t, baseURL, "/")
	assert.Equal(indexHTML, string(body), "Webroot should match expected content")
	body = fetchAndReadPage(t, baseURL, "/testfolder/test.txt")
	assert.Equal(testTXT, string(body), "Testfile should match expected content")
}

func fetchAndReadPage(t *testing.T, baseURL string, target string) string {
	require := require.New(t)

	res, err := http.Get(baseURL + target)
	require.NoError(err, "Should get page %s", target)
	require.Equal(http.StatusOK, res.StatusCode, "Should get %s", target)
	body, err := io.ReadAll(res.Body)
	require.NoError(err, "Should read page body for %s", target)
	return string(body)
}

func TestWebpageInContainerIndex(t *testing.T) {
	port := 2081
	containerName := "test-simple-fileserver-webpage-in-container-index"

	t.Parallel()
	require := require.New(t)

	err := utils.ExecCRI("run", "-d", "--rm", "-p", fmt.Sprintf("%d:8080", port), "--name", containerName, "localhost/simple-fileserver:test-e2e-webpage", "-log")
	require.NoError(err, "Failed to run the container")
	t.Cleanup(func() {
		cleanupContainer(t, containerName)
	})

	verifyPages(t, port)

	res, err := http.Get(fmt.Sprintf("http://localhost:%d/testfolder/", port))
	require.NoError(err, "Should receive response")
	require.Equal(http.StatusOK, res.StatusCode, "Should get testfolder index")
}

func TestWebpageInContainerNoIndex(t *testing.T) {
	port := 2082
	containerName := "test-simple-fileserver-webpage-in-container-no-index"

	t.Parallel()
	require := require.New(t)

	err := utils.ExecCRI("run", "-d", "--rm", "-p", fmt.Sprintf("%d:8080", port), "--name", containerName, "localhost/simple-fileserver:test-e2e-webpage", "-log", "-no-index")
	require.NoError(err, "Failed to run the container")
	t.Cleanup(func() {
		cleanupContainer(t, containerName)
	})

	verifyPages(t, port)

	res, err := http.Get(fmt.Sprintf("http://localhost:%d/testfolder/", port))
	require.NoError(err, "Should receive response")
	require.Equal(http.StatusNotFound, res.StatusCode, "Should not get testfolder index")
}

func TestWebpageAsVolumeIndex(t *testing.T) {
	port := 2083
	containerName := "test-simple-fileserver-webpage-as-volume-index"

	t.Parallel()
	require := require.New(t)

	err := utils.ExecCRI("run", "-d", "--rm", "-p", fmt.Sprintf("%d:8080", port), "--name", containerName, "-v", "./tests/e2e/testdata:/webroot:z", "localhost/simple-fileserver:test-e2e", "-log")
	require.NoError(err, "Failed to run the container")
	t.Cleanup(func() {
		cleanupContainer(t, containerName)
	})

	verifyPages(t, port)

	res, err := http.Get(fmt.Sprintf("http://localhost:%d/testfolder/", port))
	require.NoError(err, "Should receive response")
	require.Equal(http.StatusOK, res.StatusCode, "Should get testfolder index")
}

func TestWebpageAsVolumeNoIndex(t *testing.T) {
	port := 2084
	containerName := "test-simple-fileserver-webpage-as-volume-no-index"

	t.Parallel()
	require := require.New(t)

	err := utils.ExecCRI("run", "-d", "--rm", "-p", fmt.Sprintf("%d:8080", port), "--name", containerName, "-v", "./tests/e2e/testdata:/webroot:z", "localhost/simple-fileserver:test-e2e", "-log", "-no-index")
	require.NoError(err, "Failed to run the container")
	t.Cleanup(func() {
		cleanupContainer(t, containerName)
	})

	verifyPages(t, port)

	res, err := http.Get(fmt.Sprintf("http://localhost:%d/testfolder/", port))
	require.NoError(err, "Should receive response")
	require.Equal(http.StatusNotFound, res.StatusCode, "Should not get testfolder index")
}

func cleanupContainer(t *testing.T, containerName string) {
	out, err := utils.GetCommand("logs", containerName).CombinedOutput()
	if err != nil {
		t.Errorf("Failed to get logs from container %s: %v", containerName, err)
	}
	fmt.Printf("Logs from container %s:\n%s\n", containerName, out)

	err = utils.ExecCRI("rm", "-f", containerName)
	require.NoError(t, err, "Failed to remove container %s", containerName)
}

const indexHTML = `<!DOCTYPE html>
<html>
	<body>
		This is a test page.
	</body>
</html>
`
const testTXT = `This is a test file
`
