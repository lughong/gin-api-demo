package version

var (
	gitTag       = ""
	gitCommit    = "$Format:%H$" // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState = ""            // state of git tree, either "clean" or "dirty"
	buildDate    = ""            // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)
