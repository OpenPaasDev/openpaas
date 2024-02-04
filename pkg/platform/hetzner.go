package platform

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"golang.org/x/net/context"
)

type Hetzner struct{}

func (s *Hetzner) Cleanup(context.Context, *conf.Config) error { return nil }

func (s *Hetzner) Preparation(ctx context.Context, conf *conf.Config) (*conf.Config, error) {
	keysInHetzner, err := fetchHetznerKeys()
	if err != nil {
		return nil, err
	}
	keysToErase, keysToUpload := findChangesToMake(keysInHetzner, conf.CloudProviderConfig.GithubKeys)

	// erase unnecessary keys
	for _, keyId := range keysToErase {
		err := eraseKeyFromHetzner(keyId)
		if err != nil {
			fmt.Println("Error erasing key", keyId, "from Hetzner:", err)
			return nil, err
		}
	}

	// upload keys
	for _, key := range keysToUpload {
		// hetzner requires a name for a key. We don't have a name associated, so we take the last 10 chars of the key as the name
		parts := strings.SplitN(key.PublicKey, " ", 2)
		name := parts[1]
		nLen := len(name)
		if nLen > 10 {
			name = name[nLen-10:]
		}

		err := uploadKeyToHetzner(key, name)
		if err != nil {
			fmt.Println("Error uploading key", key, "to Hetzner:", err)
			return nil, err
		}
	}

	return nil, nil
}

func findChangesToMake(keysInHetzner []HetznerSSHKey, githubKeys []conf.GithubKey) ([]string, []conf.GithubKey) {
	// find any keys in github not in hetzner
	var keysToUpload []conf.GithubKey
	for _, ghKey := range githubKeys {
		if !isInHetzner(keysInHetzner, ghKey.Fingerprint) {
			keysToUpload = append(keysToUpload, ghKey)
		}
	}

	// find keys in hetzner not in github
	var keysToErase []string
	for _, k := range keysInHetzner {
		if !isInGithub(githubKeys, k.Fingerprint) {
			keysToErase = append(keysToErase, strconv.Itoa(int(k.ID)))
		}
	}

	return keysToErase, keysToUpload
}

type HetznerSSHKey struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

func fetchHetznerKeys() ([]HetznerSSHKey, error) {
	// Execute the Hetzner CLI command to list SSH keys
	out, err := exec.Command("hcloud", "ssh-key", "list", "--output", "json").Output()
	if err != nil {
		return nil, err
	}

	// Parse the output to extract key information
	var keys []HetznerSSHKey
	err = json.Unmarshal(out, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func eraseKeyFromHetzner(keyId string) error {
	// #nosec G204 we know the input we are sending to this command
	_, err := exec.Command("hcloud", "ssh-key", "delete", keyId).Output()
	return err
}

func uploadKeyToHetzner(key conf.GithubKey, name string) error {
	// #nosec G204 we know the input we are sending to this command
	_, err := exec.Command("hcloud", "ssh-key", "create", "--name", name, "--public-key", key.PublicKey).Output()
	return err
}

func isInGithub(list []conf.GithubKey, fingerprint string) bool {
	for _, v := range list {
		if v.Fingerprint == fingerprint {
			return true
		}
	}
	return false
}

func isInHetzner(list []HetznerSSHKey, fingerprint string) bool {
	for _, v := range list {
		if v.Fingerprint == fingerprint {
			return true
		}
	}
	return false
}
