package conf

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"strings"
)

func FetchGitHubKeys(githubUser string) ([]string, error) {
	url := "https://github.com/" + githubUser + ".keys"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("received non-200 status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(body), "\n"), nil
}

func GenerateFingerprint(key string) (string, error) {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key))
	if err != nil {
		return "", err
	}

	hash := md5.Sum(publicKey.Marshal())
	fingerprint := insertColons(hex.EncodeToString(hash[:]))
	return fingerprint, nil
}

func insertColons(s string) string {
	var sb strings.Builder
	for i := range s {
		if i > 0 && i%2 == 0 {
			sb.WriteByte(':')
		}
		sb.WriteByte(s[i])
	}
	return sb.String()
}
