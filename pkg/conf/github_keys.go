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

func UpdateConfigWithGithubKeys(cnf *Config) (*Config, error) {
	// get the github ids in the config, retrieve the keys
	for _, user := range cnf.CloudProviderConfig.GithubIds {
		keys, err := fetchGitHubKeys(user)
		if err != nil {
			return nil, err
		}
		// convert the keys to fingerprints
		var fingerprints []string
		for _, key := range keys {
			fprint, err := generateFingerprint(key)
			if err != nil {
				return nil, err
			}
			fingerprints = append(fingerprints, fprint)
		}
		// add the fingerprints to the config
		if sshKeys, ok := cnf.CloudProviderConfig.ProviderSettings["ssh_keys"]; ok {
			// checks the value in the map is a slice of string
			switch keys := sshKeys.(type) {
			case []string:
				cnf.CloudProviderConfig.ProviderSettings["ssh_keys"] = append(keys, fingerprints...)
			default:
				// Handle the case where ssh_keys exists but is not a slice of strings, which would be wrong
				cnf.CloudProviderConfig.ProviderSettings["ssh_keys"] = fingerprints
			}
		} else {
			// Create ssh_keys with the new keys if it doesn't exist
			cnf.CloudProviderConfig.ProviderSettings["ssh_keys"] = fingerprints
		}
	}
	return cnf, nil
}

func fetchGitHubKeys(githubUser string) ([]string, error) {
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

func generateFingerprint(key string) (string, error) {
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
