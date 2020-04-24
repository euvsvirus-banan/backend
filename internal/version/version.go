package version

// Variables set using ldflags to the go build command-line.
var (
	Project     string // nolint: gochecknoglobals
	Version     string // nolint: gochecknoglobals
	GitRevision string // nolint: gochecknoglobals
	BuildDate   string // nolint: gochecknoglobals
	GoVersion   string // nolint: gochecknoglobals
)
