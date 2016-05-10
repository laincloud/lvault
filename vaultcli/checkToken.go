package vaultcli

import (
	"time"

	"github.com/mijia/sweb/log"
)

//检查是否是 root token
func (c *VaultClient) CheckRootToken(token string) bool {
	return c.CheckToken(token, "root")
}

// 前置条件必须是 vault 集群至少有一个节点是解锁的
func (c *VaultClient) CheckToken(token string, tokentype string) bool {
	for i := 0; i < 5; i = i + 1 {
		if c.checkToken(token, tokentype) {
			return true
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return false
}

func (c *VaultClient) checkToken(token string, tokentype string) bool {
	c.c.SetToken(token)
	defer c.c.ClearToken()
	tokenAuth := c.c.Auth().Token()
	secret, err := tokenAuth.LookupSelf()
	if err != nil {
		log.Error(err)
		return false
	}
	if tokentype == "root" {
		if secret.LeaseDuration != 0 {
			return false
		}
	} else {

	}
	return true
}
