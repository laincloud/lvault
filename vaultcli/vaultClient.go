package vaultcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/api"
)

type VaultClient struct {
	status     VaultStatus
	c          *api.Client
	httpClient *http.Client
}

// RawList list secrets in vault
// We reimplement it because VaultClient.c will wait for 10 seconds
func (c *VaultClient) RawList(token, path string) (*api.Secret, error) {
	req, err := http.NewRequest("LIST", fmt.Sprintf("%sv1/%s", VaultURL, path), nil)
	if err != nil {
		return nil, err
	}

	setVaultToken(req, token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return api.ParseSecret(resp.Body)
}

// RawRead read secret in vault
// We reimplement it because VaultClient.c will wait for 10 seconds
func (c *VaultClient) RawRead(token, path string) (*api.Secret, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sv1/%s", VaultURL, path), nil)
	if err != nil {
		return nil, err
	}

	setVaultToken(req, token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return api.ParseSecret(resp.Body)
}

// RawUnseal unseal vault
// We reimplement it because VaultClient.c will wait for 10 seconds
func (c *VaultClient) RawUnseal(shard string) (*api.SealStatusResponse, error) {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"key": shard,
	}
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%sv1/sys/unseal", VaultURL), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result api.SealStatusResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func setVaultToken(r *http.Request, token string) {
	r.Header.Set("X-Vault-Token", token)
}
