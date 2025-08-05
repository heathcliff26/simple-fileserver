package e2e

import (
	"fmt"
	"os"
	"testing"

	"github.com/heathcliff26/simple-fileserver/tests/utils"
)

func TestMain(m *testing.M) {
	err := os.Chdir("../..")
	if err != nil {
		fmt.Println("Failed to change directory to project root: %w", err)
		os.Exit(1)
	}

	err = utils.ExecCRI("build", "-t", "localhost/simple-fileserver:test-e2e", ".")
	if err != nil {
		fmt.Println("Failed to build the image: %w", err)
		os.Exit(1)
	}

	err = utils.ExecCRI("build", "-t", "localhost/simple-fileserver:test-e2e-webpage", "tests/e2e/testdata")
	if err != nil {
		fmt.Println("Failed to build the image: %w", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
