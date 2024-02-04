package conf

import (
	"crypto/md5" // #nosec G501 we need md5 as this is how hetzner calculates fingerprints
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"golang.org/x/crypto/ssh"
)

func UpdateConfigWithGithubKeys(ctx context.Context, cnf *Config) (*Config, error) {
	// get the github ids in the config, retrieve the keys
	for _, user := range cnf.CloudProviderConfig.GithubIds {
		keys, err := fetchGitHubKeys(ctx, user)
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

func fetchGitHubKeys(ctx context.Context, githubUser string) ([]string, error) {
	const baseURL = "https://github.com"
	fullURL := fmt.Sprintf("%s/%s.keys", baseURL, url.PathEscape(githubUser))
	// #nosec G107 we are loading a dynamic url in purpose
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// we need to close the body to avoid issues, but linter complains if we ignore the error case
	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

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

	// #nosec G401 we need md5 as this is how hetzner calculates fingerprints
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
