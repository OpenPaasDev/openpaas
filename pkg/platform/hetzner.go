package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/OpenPaasDev/openpaas/pkg/util"

	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"golang.org/x/net/context"
)

const GHKeyPrefix string = "gh-key-"

type Hetzner struct{}

func (s *Hetzner) Cleanup(context.Context, *conf.Config) error {
	fmt.Println("Cleanup for Hetzner platform...")
	return nil
}

func (s *Hetzner) Prepare(ctx context.Context, conf *conf.Config) error {
	fmt.Println("Preparing Hetzner platform...")
	err := runPreparationLogic(ctx, conf, fetchHetznerKeys, fetchGitHubKeys, eraseKeyFromHetzner, uploadKeyToHetzner, util.GetPublicIP)
	if err != nil {
		return err
	}
	return nil
}

type HetznerSSHKey struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

func runPreparationLogic(ctx context.Context,
	conf *conf.Config,
	getHetznerKeys func() ([]HetznerSSHKey, error),
	getGithubKeys func(context.Context, string) ([]string, error),
	eraseHetznerKey func(HetznerSSHKey) error,
	uploadHetznerKey func(string, string) error,
	getPublicIp func(context.Context) (string, error),
) error {
	err := setHetznerContext(conf.CloudProviderConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("Error setting Hetzner context: %v\n", err))
	}
	keysInHetzner, err := getHetznerKeys()
	if err != nil {
		return errors.New(fmt.Sprintf("Error fetching keys from Hetzner: %v\n", err))
	}

	// load keys from github
	var githubPublicKeys []string
	for _, user := range conf.CloudProviderConfig.GithubIds {
		keys, er := getGithubKeys(ctx, user)
		if er != nil {
			return er
		}
		githubPublicKeys = append(githubPublicKeys, keys...)
	}

	// delete keys that start with the github prefix, as they could be outdated from a previous upload
	for _, key := range keysInHetzner {
		if strings.HasPrefix(key.Name, GHKeyPrefix) {
			er := eraseHetznerKey(key)
			if er != nil {
				fmt.Println("Error erasing key", key.ID, "with name", key.Name, "from Hetzner:", er)
				return er
			}
		}
	}

	// upload all github keys to hetzner
	for _, publicKey := range githubPublicKeys {
		// hetzner requires a name for a key. We don't have a name associated, so we take the last 10 chars of the key as the name
		parts := strings.SplitN(publicKey, " ", 2)
		name := parts[1]
		nLen := len(name)
		if nLen > 10 {
			name = GHKeyPrefix + name[nLen-10:]
		}
		er := uploadHetznerKey(publicKey, name)
		if er != nil {
			fmt.Println("Error uploading key ending in", name, "to Hetzner:", er)
			return er
		}
	}

	// read the keys again to extract the ids assigned to the new keys
	updatedKeysInHetzner, err := getHetznerKeys()
	if err != nil {
		return err
	}
	var githubIds []string
	for _, key := range updatedKeysInHetzner {
		if strings.HasPrefix(key.Name, GHKeyPrefix) {
			githubIds = append(githubIds, strconv.Itoa(int(key.ID)))
		}
	}

	// update the config adding the github key ids to any existing config in place
	conf.CloudProviderConfig.ProviderSettings["ssh_keys"] = append([]string{}, githubIds...)

	// update allowed ips record with the current public ip
	ip, err := getPublicIp(ctx)
	if err != nil {
		return err
	}
	conf.CloudProviderConfig.AllowedIPs = append(conf.CloudProviderConfig.AllowedIPs, ip+"/32")

	return nil
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
func setHetznerContext(config conf.CloudProvider) error {
	hetznerCtx, ok := config.ProviderSettings["context"].(string)
	if !ok {
		fmt.Println("Warning: no hetzner context found in provider settings, using currently set context")
		return nil
	}
	_, err := exec.Command("hcloud", "context", "use", hetznerCtx).Output() //nolint
	return err
}

func eraseKeyFromHetzner(key HetznerSSHKey) error {
	fmt.Println("Erasing from Hetzner key", key.ID, "with name", key.Name)
	// #nosec G204 we know the input we are sending to this command
	_, err := exec.Command("hcloud", "ssh-key", "delete", strconv.Itoa(int(key.ID))).Output()
	return err
}

func uploadKeyToHetzner(publicKey string, name string) error {
	fmt.Println("Uploading to Hetzner public key ending in", name)
	// #nosec G204 we know the input we are sending to this command
	_, err := exec.Command("hcloud", "ssh-key", "create", "--name", name, "--public-key", publicKey).Output()
	return err
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
