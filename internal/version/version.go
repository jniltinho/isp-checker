package version

import (
	"fmt"
	"io"
	"runtime"

	"isp_checker"
)

// PrintVersion prints the version details to a specified writer
func PrintVersion(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Version: %s\n", isp_checker.BuildVersion)
	_, _ = fmt.Fprintf(w, "Git Commit Hash: %s\n", isp_checker.BuildHash)
	_, _ = fmt.Fprintf(w, "Build Date: %s\n", isp_checker.BuildDate)
	_, _ = fmt.Fprintf(w, "OS: %s\n", runtime.GOOS)
	_, _ = fmt.Fprintf(w, "Architecture: %s\n", runtime.GOARCH)
}
