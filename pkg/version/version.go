package version

import (
	"runtime"
	"runtime/debug"
	"strings"
)

const Name = "simple-fileserver"

// NOTE: The $Format strings are replaced during 'git archive' thanks to the
// companion .gitattributes file containing 'export-subst' in this same
// directory.  See also https://git-scm.com/docs/gitattributes
var gitCommit string = "$Format:%H$" // sha1 from git, output of $(git rev-parse HEAD)
var gitVersion string = ""

func init() {
	initGitCommit()
	initGitVersion()
}

func initGitCommit() {
	if strings.HasPrefix(gitCommit, "$Format") {
		var commit string
		buildinfo, _ := debug.ReadBuildInfo()
		for _, item := range buildinfo.Settings {
			if item.Key == "vcs.revision" {
				commit = item.Value
				break
			}
		}
		if commit == "" {
			commit = "Unknown"
		}
		gitCommit = commit
	}
	if gitCommit == "" {
		gitCommit = "Unknown"
	}
}

func initGitVersion() {
	if gitVersion == "" {
		buildinfo, _ := debug.ReadBuildInfo()
		gitVersion = buildinfo.Main.Version
	}
}

// Return a formatted string containing the version, git commit and go version the app was compiled with.
func Version() string {
	commit := gitCommit
	if len(commit) > 7 {
		commit = commit[:7]
	}

	result := Name + ":\n"
	result += "    Version: " + gitVersion + "\n"
	result += "    Commit:  " + commit + "\n"
	result += "    Go:      " + runtime.Version() + "\n"

	return result
}
