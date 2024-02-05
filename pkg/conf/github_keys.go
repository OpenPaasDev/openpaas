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

type GithubKey struct {
	PublicKey   string
	Fingerprint string
}

func UpdateConfigWithGithubKeys(ctx context.Context, cnf *Config) (*Config, error) {
	var githubKeys []GithubKey
	var fingerprints []string

	// get the github ids in the config, retrieve the keys
	for _, user := range cnf.CloudProviderConfig.GithubIds {
		keys, err := fetchGitHubKeys(ctx, user)
		if err != nil {
			return nil, err
		}

		// convert the keys to fingerprints and store the data
		for _, key := range keys {
			fprint, err := generateFingerprint(key)
			if err != nil {
				return nil, err
			}
			fingerprints = append(fingerprints, fprint)
			githubKeys = append(githubKeys, GithubKey{key, fprint})
		}
	}
	// add the github keys we obtained, in case we need to upload them to the provider
	cnf.CloudProviderConfig.GithubKeys = githubKeys

	// the yaml loaded entry is a []interface{} not []string{} so we need to convert it, if present
	var sshKeysFromConfig []string
	if sshKeys, ok := cnf.CloudProviderConfig.ProviderSettings["ssh_keys"]; ok {
		sshKeysFromConfig = interfaceToStringSlice(sshKeys.([]interface{}))
	}
	// append fingerprints to provided ssh_keys
	cnf.CloudProviderConfig.ProviderSettings["ssh_keys"] = append(sshKeysFromConfig, fingerprints...)

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

	keys := strings.Split(string(body), "\n")
	// Filter out empty strings
	var nonEmptyKeys []string
	for _, part := range keys {
		if part != "" {
			nonEmptyKeys = append(nonEmptyKeys, part)
		}
	}
	return nonEmptyKeys, nil
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

func interfaceToStringSlice(in []interface{}) []string {
	var slice []string
	for _, key := range in {
		strKey, ok := key.(string)
		if ok {
			slice = append(slice, strKey)
		}
	}
	return slice
}
