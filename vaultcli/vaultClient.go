package vaultcli

import "github.com/hashicorp/vault/api"

type VaultClient struct {
	status VaultStatus
	c      *api.Client
}
