package util

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type IP struct { //nolint
	Query string `json:"query"`
}

func init() {
	rand.Seed(time.Now().UnixNano()) //nolint
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))] //nolint
	}
	return string(b)
}

func GetPublicIP(ctx context.Context) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://ip-api.com/json/", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			panic(e)
		}
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	var ip IP
	err = json.Unmarshal(body, &ip)
	if err != nil {
		return "", err
	}
	// fmt.Print(ip.Query)
	return ip.Query, nil
}

func RunCmd(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fmt.Println("Error running:", cmd, args, err)
		return err
	}
	return nil
}

func IsBrewInstalled() bool {
	err := RunCmd("brew", "--version")
	if err != nil {
		return false
	}
	return true
}

func IsPipInstalled() bool {
	err := RunCmd("pip", "--version")
	if err != nil {
		return false
	}
	return true
}

func IsAnsibleInstalled() bool {
	err := RunCmd("ansible", "--version")
	if err != nil {
		return false
	}
	return true
}

func IsHCloudInstalled() bool {
	err := RunCmd("hcloud", "version")
	if err != nil {
		return false
	}
	return true
}

func ExtractTarGz(tempPath string, destPath string) error {
	// Limit the size of reader output to 100 Mb to prevent decompression bombs
	limitForReaders := int64(100 * 1024 * 1024)
	file, err := os.Open(filepath.Clean(tempPath))
	if err != nil {
		fmt.Println("Error opening downloaded file:", err)
		return err
	}
	defer file.Close() //nolint

	limiter := io.LimitReader(file, limitForReaders)
	uncompressedStream, err := gzip.NewReader(limiter)
	if err != nil {
		return err
	}

	tarLimiter := io.LimitReader(uncompressedStream, limitForReaders)
	tarReader := tar.NewReader(tarLimiter)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Construct the full path and create the file
		path := filepath.Join(destPath, filepath.Clean(header.Name))
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0750); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Clean(path))
			if err != nil {
				return err
			}
			tarReaderLimiter := io.LimitReader(tarReader, limitForReaders)
			if _, err := io.Copy(outFile, tarReaderLimiter); err != nil {
				outFile.Close() // nolint
				return err
			}
			outFile.Close() // nolint
		}
	}
	return nil
}

func DownloadFile(ctx context.Context, url, outputPath string) error {
	// Get the data
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil) // #nosec G107 called internally with known urls
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return err
	}
	defer resp.Body.Close() // nolint

	// Create the file
	out, err := os.Create(filepath.Clean(outputPath))
	if err != nil {
		return err
	}
	defer out.Close() // nolint

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
