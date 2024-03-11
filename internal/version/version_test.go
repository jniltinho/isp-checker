package version_test

import (
	"bufio"
	"io"
	"testing"

	"isp_checker/internal/version"

	"github.com/stretchr/testify/assert"
)

func TestPrintVersion(t *testing.T) {
	wr := bufio.NewWriter(io.Discard)
	version.PrintVersion(wr)
	assert.Equal(t, 71, wr.Buffered())
	assert.NoError(t, wr.Flush())
}
