package vaultcli

import (
	"container/list"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mijia/sweb/log"
)

var VaultURL string

func (c *VaultClient) InitClient(tls bool) {
	VaultURL = "http://lvault." + os.Getenv("LAIN_DOMAIN") + "/"
	if tls {
		log.Info("using https to send messages to vault cluster")
		VaultURL = "https://lvault." + os.Getenv("LAIN_DOMAIN") + "/"
	}
	config := api.DefaultConfig()
	config.Address = VaultURL
	var err error
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	c.c = client
	go c.status.UpdateStatus()
	go c.UpdateClient()
}

func (c *VaultClient) UpdateClient() {
	for {
		time.Sleep(5 * time.Second)
		seal, err := c.c.Sys().SealStatus()
		if err != nil {
			if strings.Contains(err.Error(), "server is not yet initialized") {
				continue
			}
			panic(err)
		}
		if seal.Sealed == false {
			log.Debug("healthy")
			continue
		}
		unsealedURL := c.status.UnsealedURL()
		if unsealedURL == "" {
			log.Debug("find empty")
			continue
		}
		config := api.DefaultConfig()
		config.Address = unsealedURL
		log.Info("change client url to " + unsealedURL)
		client, err := api.NewClient(config)
		if err != nil {
			panic(err)
		}
		c.c = client
	}
}

func (c *VaultClient) InitVault(req *api.InitRequest) (*api.InitResponse, error) {
	return c.c.Sys().Init(req)
}

func (c *VaultClient) Unseal(unsealkey []string) (err error) {
	for {
		if c.status.AllUnsealed() {
			break
		}
		for _, v := range unsealkey {
			err = c.unseal(v)
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

func (c *VaultClient) unseal(unsealkey string) error {
	_, err := c.c.Sys().Unseal(unsealkey)
	return err
}

func (c *VaultClient) PutSecret(token string, path string, data string) error {
	c.c.SetToken(token)
	defer c.c.ClearToken()
	secrets := make(map[string]interface{})
	secrets["value"] = data
	_, err := c.c.Logical().Write(path, secrets)
	return err
}

func (c *VaultClient) DeleteSecret(token string, path string) error {
	c.c.SetToken(token)
	defer c.c.ClearToken()
	_, err := c.c.Logical().Delete(path)
	return err
}

func (c *VaultClient) ListSecrets(token string, path string) ([]string, error) {
	c.c.SetToken(token)
	defer c.c.ClearToken()
	logical := c.c.Logical()
	l := list.New()
	l.PushBack(path)
	ret := []string{}
	for l.Len() > 0 {
		iter := l.Back()
		p := iter.Value.(string)
		l.Remove(iter)
		s, err := logical.List(p)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if s == nil {
			log.Debug(p)
		}
		if s == nil || len(s.Data) == 0 {
			s, err = logical.Read(p)
			log.Debug(s)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			if s != nil {
				if len(s.Data) != 0 {
					ret = append(ret, s.Data["value"].(string))
				}
			}
		} else {
			log.Debug(s.Data)
			keys := s.Data["keys"].([]interface{})
			for _, v := range keys {
				var newpath string
				if strings.HasSuffix(p, "/") {
					newpath = p + v.(string)
				} else {
					newpath = p + "/" + v.(string)
				}
				l.PushBack(newpath)
			}
		}
	}
	return ret, nil
}
