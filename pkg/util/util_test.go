package util

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandString(t *testing.T) {

	theMap := make(map[string]string)

	for i := 0; i <= 30; i++ {
		str := RandString(20)
		if _, ok := theMap[str]; ok {
			assert.True(t, false)
		} else {
			theMap[str] = str
		}
	}
}

func TestGetPublicIP(t *testing.T) {
	ip, err := GetPublicIP(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 3, strings.Count(ip, "."))
}

func Test_RunCmd(t *testing.T) {
	err := RunCmd("ls", "-l")
	require.NoError(t, err)
}

func Test_ExtractTarGZ(t *testing.T) {
	tempPath := "../testdata/testfile.tar.gz"
	destPath := "/tmp"
	err := ExtractTarGz(tempPath, destPath)
	require.NoError(t, err)
	content, err := os.ReadFile(destPath + "/test.txt")
	require.NoError(t, err)
	text := string(content)
	assert.Equal(t, "test", text)
}

func Test_DownloadFile(t *testing.T) {
	ctx := context.Background()
	url := "https://httpbin.org/robots.txt"
	destPath := "/tmp/test-download.txt"
	expectedContent := "User-agent: *\nDisallow: /deny\n"

	err := DownloadFile(ctx, url, destPath)
	require.NoError(t, err)
	content, err := os.ReadFile(destPath)
	require.NoError(t, err)
	text := string(content)
	assert.Equal(t, expectedContent, text)
}
